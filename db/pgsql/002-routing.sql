-- Copyright © 2014–5 Brad Ackerman.
--
-- Licensed under the Apache License, Version 2.0 (the "License");
-- you may not use this file except in compliance with the License.
-- You may obtain a copy of the License at
--
-- http://www.apache.org/licenses/LICENSE-2.0
--
-- Unless required by applicable law or agreed to in writing, software
-- distributed under the License is distributed on an "AS IS" BASIS,
-- WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
-- See the License for the specific language governing permissions and
-- limitations under the License.

-- Functions for wrapping pgRouting.
BEGIN;
-- solarsystem_route_map converts solar system IDs from the CCP-assigned
-- values (globally-defined) to the PGRouting topology vertex IDs (random).
DROP MATERIALIZED VIEW IF EXISTS solarsystem_route_map;
CREATE MATERIALIZED VIEW solarsystem_route_map AS
  SELECT mss."solarSystemID" ccpid, v.id pgrid
  FROM "mapSolarSystems" mss, "mapSolarSystemJumps_vertices_pgr" v
  WHERE mss.the_geom && v.the_geom;
CREATE INDEX ON solarsystem_route_map USING hash (ccpid);
CREATE INDEX ON solarsystem_route_map USING hash (pgrid);

-- Find the route from two system IDs.
CREATE OR REPLACE FUNCTION eve_findRoute(
  IN srcSystemID integer,
  IN destSystemID integer,
  OUT seq integer,
  OUT systemID integer,
  OUT systemName text
) RETURNS SETOF RECORD AS
$$
DECLARE
  rec record;
  sql text;
  source integer;
  target integer;
BEGIN
  EXECUTE 'SELECT m1.pgrid src, m2.pgrid dst ' ||
          'FROM   solarsystem_route_map m1, solarsystem_route_map m2 ' ||
          'WHERE  m1.ccpid = ' || srcSystemID ||
          ' AND   m2.ccpid = ' || destSystemID
          INTO rec;
  source := rec.src;
  target := rec.dst;
  IF source IS NULL OR target IS NULL then
    -- one of these systems isn't on the jumpgate network.
    RETURN;
  END IF;
  seq := 0;
  sql := 'WITH route AS (SELECT seq, id1 AS node FROM pgr_dijkstra(' ||
         -- The SELECT statement for the jumps is passed into pgr_dijkstra
         -- as text.
  			 '''SELECT id, source, target, 1 :: float8 AS cost, x1, y1, x2, y2 ' ||
  			 'FROM "mapSolarSystemJumps"'', ' ||
         source || ', ' || target ||
         ', true, false)) ' ||
         'SELECT r.seq, mss."solarSystemID" id, mss."solarSystemName" ssName ' ||
         'FROM route r, "mapSolarSystems" mss, solarsystem_route_map rm ' ||
         'WHERE mss."solarSystemID" = rm.ccpid AND rm.pgrid = r.node';
  FOR rec IN EXECUTE sql
  LOOP
    seq := seq + 1;
    systemID := rec.id;
    systemName := rec.ssName;
    RETURN NEXT;
  END LOOP;
  RETURN;
END;
$$ LANGUAGE plpgsql IMMUTABLE;

COMMIT;

-- Vacuum now that we've indexed; can't run in a transaction block.
VACUUM ANALYZE solarsystem_route_map;

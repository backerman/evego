-- Copyright © 2014–6 Brad Ackerman.
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

BEGIN;

-- solarsystem_route_map interferes with redoing the topology, so we
-- drop at the beginning and add back afterwards.
DROP MATERIALIZED VIEW IF EXISTS solarsystem_route_map;

UPDATE "mapSolarSystems" SET the_geom=ST_MakePoint(x,y);

UPDATE ONLY "mapSolarSystemJumps" j
  SET x1 = fss.x, y1 = fss.y, x2 = tss.x, y2 = tss.y,
       source = j."fromSolarSystemID", target = j."toSolarSystemID",
       the_geom = ST_MakeLine(fss.the_geom, tss.the_geom)
  FROM "mapSolarSystems" fss, "mapSolarSystems" tss
  WHERE j."fromSolarSystemID" = fss."solarSystemID"
  AND   j."toSolarSystemID" = tss."solarSystemID";

-- Create topology
-- units in meters and interstellar distances, so tolerance is a tad larger
-- than pgrouting's usual applications.
SELECT pgr_createTopology('mapSolarSystemJumps', 1e6);

CREATE MATERIALIZED VIEW solarsystem_route_map AS
  SELECT mss."solarSystemID" ccpid, v.id pgrid
  FROM "mapSolarSystems" mss, "mapSolarSystemJumps_vertices_pgr" v
  WHERE mss.the_geom && v.the_geom;
CREATE INDEX ON solarsystem_route_map USING hash (ccpid);
CREATE INDEX ON solarsystem_route_map USING hash (pgrid);

COMMIT;

-- Vacuum now that we've indexed; can't run in a transaction block.
VACUUM ANALYZE solarsystem_route_map;

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

-- Enable spatial functionality in database
CREATE EXTENSION IF NOT EXISTS postgis;
CREATE EXTENSION IF NOT EXISTS pgrouting;

BEGIN;

-- Create geometry column
-- No Z here because pgrouting doesn't support 3D yet.
ALTER TABLE "mapSolarSystems" ADD COLUMN the_geom geometry;
UPDATE "mapSolarSystems" SET the_geom=ST_MakePoint(x,y);

-- Index the geometry column we just added.
CREATE INDEX ON "mapSolarSystems" USING GIST (the_geom);

ALTER TABLE "mapSolarSystemJumps"
  ADD COLUMN id serial,
  ADD COLUMN cost double precision default 1.0,
  ADD COLUMN x1 double precision,
  ADD COLUMN y1 double precision,
  ADD COLUMN x2 double precision,
  ADD COLUMN y2 double precision,
  ADD COLUMN source int4,
  ADD COLUMN target int4,
  ADD COLUMN the_geom geometry;

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
COMMIT;

-- Vacuum now that we've indexed; can't run in a transaction block.
VACUUM ANALYZE "mapSolarSystems";

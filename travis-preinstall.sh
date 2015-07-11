#!/bin/sh
#
# Copyright © 2014–5 Brad Ackerman.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
# ---
#
# Environment setup script for Travis CI. This is moderately evil,
# but Travis is running Ubuntu 12LTS and that's got too-old sqlite3
# so nothing will work without such steps.

sudo add-apt-repository -y ppa:ubuntugis/ppa &&
sudo apt-get update -qq &&
sudo apt-get install -y libproj-dev make libxml2-dev zlib1g-dev \
  pkg-config libgeos-c1 libgeos-dev &&
wget http://www.sqlite.org/2015/sqlite-autoconf-3080801.tar.gz &&
tar -zxvf sqlite-autoconf-3080801.tar.gz &&
(cd sqlite-autoconf-3080801 && ./configure && make && sudo make install) &&
sudo cp /usr/local/lib/libsqlite3.so.0.8.6 /usr/lib/x86_64-linux-gnu &&
wget http://www.gaia-gis.it/gaia-sins/libspatialite-sources/libspatialite-4.3.0.tar.gz &&
tar -zxvf libspatialite-4.3.0.tar.gz &&
cd libspatialite-4.3.0 &&
./configure --enable-geosadvanced=no --enable-freexl=no && make &&
sudo make install

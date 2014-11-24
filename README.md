# evego

evego is a library for your [Internet spreadsheets][eve] spreadsheets;
sometimes, a spreadsheet isn't actually sufficient. It does not currently slice,
dice, or make Julienne fries. For that matter, it doesn't even interface with
the user; it's just a back-end library.

[eve]: http://www.eveonline.com/

## Development

We like test cases. Unit tests are written using [GoConvey][convey], and there
should be as close to 100% coverage as possible. Higher-level tests would also
be a good idea.

[convey]: http://goconvey.co/

This repository includes the subset of the [EVE Static Data Export][sde]
necessary to run the test cases. If you add test cases that use data not in
this subset:

* Add the missing items (and tables, if necessary) to `make-test-db.sh`;
* Rerun that script against [Fuzzysteve][steve]'s SQLite [conversion] of the
full SDE; and
* Add the new version of `testdb.sqlite` to your changeset.

[conversion]: https://www.fuzzwork.co.uk/dump/
[sde]: https://developers.eveonline.com/resource/static-data-export
[steve]: https://www.fuzzwork.co.uk/

## License

Portions of the EVE static data export are included in this repository
(`testdb.sqlite`); the following notice applies to that file:

© 2014 CCP hf. All rights reserved. "EVE", "EVE Online", "CCP", and all related
logos and images are trademarks or registered trademarks of CCP hf.

The remainder of this repository is © 2014 Brad Ackerman and licensed under the
[Apache License 2.0][apache], the full text of which is in the LICENSE file.

[apache]: http://www.apache.org/licenses/LICENSE-2.0

# Mainflux Licensing Service

[![coverage][cov-badge]][cov-url]
![build][build-badge]

Licensing service

[cov-badge]: https://codecov.io/gh/mainflux/license/branch/master/graph/badge.svg
[cov-url]: https://codecov.io/gh/mainflux/license
[build-badge]: https://github.com/mainflux/license/workflows/build/badge.svg

Mainflux license service is self-contained, general-purpose, independent service license management tool. It provides API for CRUD operations, as well as License validation endpoint. License object consists of the following fields:

- ID - unique License identifier
- Key - License key given alongside License ID used for validation
- Issuer - the ID of the issuer
- Device ID - Unique identifier of the device License is issued for
- Active - flag that indicates whether the license is active or not
- Created at - UTC date and time the when the License is issued
- Expires at - UTC date and time after which the License is not valid anymore
- Updated by - id of the user who executed the last License update
- Updated at - UTC date and time of the last License update.
- Services - list of services the License applies to
- Plan - additional License metadata in JSON format

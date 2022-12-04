# what's this?

scraping [L-Gate](https://www.info.l-gate.net/), dumps and save csv file.

**ALL FOR GIGASCHOOL, LET'S SHARING PHILOSOPHY.**

# API（WIP)

The following endpoints are confirmed to exist.


- `https://{host}-api.l-gate.net`
    - API base url
- `/control`
    - Perhaps this is the basis for all endpoints
- `/control/auth/login`
    - For login. **POST** method
- `/control/manual/get`
    - I'm not really sure, but it seems that special values (named `CONTROLSESSID`) required for Cookies are returned in the Http Response. This value is used for file requests and other endpoints
- `/control/action-log/download-csv-total`
    - All user action log data api endpoint
    - Some parameters required, for example, `?start_at=00000&end_at=999999&time_unit=hour&scope=tenant&action=&response_all=1&encoding=utf8"`
- `/control/user/download-csv`
    - All registed user data api endpoint
    - Some parameters required, almost same above example, like `?encoding=utf8&term_uuid=xxxxxx-xxxx-xxxx-xxxx-xxxxx&page_size=50`
- `/control/job-state/view/{uuid}`
    -  Seems to be asking about the readiness of the file.
    - This UUID will be included in the response when calling the API for data acquisition
- `/control/file/view/{uuid}`
    - CSV file here.
    - The last response of above `/control/job-state/view/{uuid}` has **Result** Object property, and if all process would be done, **Result UUID** is included in the property
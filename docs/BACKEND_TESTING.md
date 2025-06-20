# Backend testing

## Testing flow

Backend testing can be done by defining test files in the same folder as the test entry point. Starting the test from the entry point facilitate the database connection and testing via transaction instead of mocking the database. This provides more reliable and assuring test results, especially when the business logic is not drilled into any service layer.

Here's one the reasonable flows:
* Ensure database exist and the service is buildable
* Write migration scripts, this is for transaction
* Seed the data (this does not require an external script, just define in entrypoint)
    * Push the fixtures in the entry point. The problem is, this blows with more data.
    * Do transaction and run the codes.
    * Ensure no commit is done, or the data will be back again.
* run main test
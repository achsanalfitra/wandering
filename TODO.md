# Building the vector machine service

1. The model is ready, initiate the database so its ready.
2. ~~Define the cannonical_order maintenance endpoints in the main app: ~~the minimum is POST (create),~~ GET (print), DELETE(remove a vibe and key), and UPDATE (change status, name or order)~~
3. ~~Add multi insertion POST method~~
4. Add forwarding endpoint:
   1. ~~Validate the input data, check existence of the array given by the frontend~~
   2. ~~Return OK or NOT OK~~
   3. When the service is created
5. Test every handlers

# Key requirement tasks

1. create testing architecture within this codebase
2. set up test environment, including db and other dependencies
3. define test workflow
# BConnect-backend

# BConnect
Webapp that aims the facilitate social networking between UCLA alumni

# Setup
Assuming you have Go set up and installed and your $GOPATH environment variable exists, clone this repository in the $GOPATH/src/github.com/ucladevx folder. Install postgres on your local computer and start it up. `cd $GOPATH/src/github.com/ucladevx/BConnect-backend` and  run `/scripts/init.sh`. This should populate a database on your local computer called `connect_b_users` with entries. You can check this by running `psql postgres` and in the psql command prompt, type `\c connect_b_users connect_b`, then `\dt` to see the list of tables you have. 
To run to backend for now, run `make bin/bconnect` and `./bin/BConnect-backend` from the project root directory.

# Current Progress
Data storage with postgres
Simple authentication flow with JWT
Simple friend schema with friend requests and separate friend db
Basic filtering

# Upcoming features
More friend functionality
Chat
Location-based functionality
Use Docker 

# Expense Tracker

A Telegram Bot to track your expenses.

[![Build Status](https://travis-ci.com/masudur-rahman/expense-tracker-bot.svg?branch=main)](https://travis-ci.com/masudur-rahman/expense-tracker-bot)
[![Go Report Card](https://goreportcard.com/badge/github.com/masudur-rahman/expense-tracker-bot)](https://goreportcard.com/report/github.com/masudur-rahman/expense-tracker-bot)
[![codecov](https://codecov.io/gh/masudur-rahman/expense-tracker-bot/branch/master/graph/badge.svg)](https://codecov.io/gh/masudur-rahman/expense-tracker-bot)
[![License](https://img.shields.io/github/license/masudur-rahman/expense-tracker-bot)](


## Description

`Expense Tracker Bot` is a Telegram Bot to track your expenses. It is built using [Go](https://golang.org/) and [Postgres](https://www.postgresql.org/).

It's currently supporting a single user [masudur-rahman](https://t.me/masudur_rahman).


## Requirements

### Telegram Bot

1. Create a new bot using [BotFather](https://t.me/botfather).
   - Use `/newbot` command to create a new bot.
   - Use `/setname` command to set a name for the bot.
   - Use `/setdescription` command to set a description for the bot.

2. Set commands for the bot using `/setcommands` command.
    ```
    new - Add new Transaction, Account or User
    newtxn - Add new transaction
    user - List persons involved in some loan/borrow with the system user
    balance - List Account Balance
    expense - Fetch Expense of Current month
    summary - Transaction summary of current month
    allsummary - Transaction summary based on Type, Category, Subcategory
    report - Transaction Report
    cat - List Transaction categories
    ```
 
3. Create a Token for the bot.
    - Use `/token` command to get the bot token.

### Postgres Database

#### Local Setup

- Install Postgres using [Homebrew](https://brew.sh/).
   ```bash
   brew update
   brew install postgresql
   brew services start postgresql
   ```
- Create user superuser `postgres` with password `postgres`.
   ```bash
   psql postgres -c "CREATE USER postgres WITH SUPERUSER PASSWORD 'postgres';"
   ```
- Create a new database named `expense`.
   ```bash
   psql -u postgres -c "CREATE DATABASE expense;"
   ```

## Installation and Running

### Local Setup

1. Clone the repository.
    ```bash
    mkdir -p $GOPATH/src/github.com/masudur-rahman
    cd $GOPATH/src/github.com/masudur-rahman
    git clone git@github.com:masudur-rahman/expense-tracker-bot.git
    ```
2. Export required environment variables.
    ```bash
    export TELEGRAM_BOT_TOKEN=<TELEGRAM_BOT_TOKEN>
    
    # following environment variables are the default values
    # you can ignore them if you are using the default values
    export POSTGRES_USER=postgres
    export POSTGRES_PASSWORD=postgres
    export POSTGRES_DB=expense
    export POSTGRES_HOST=localhost
    export POSTGRES_PORT=5432
    export POSTGRES_SSL_MODE=disable
    ```
3. Run `make run` to start the server.

### Back4App Setup

1. Create a new app on [Back4App](https://www.back4app.com/).
2. Connect your app to a GitHub repository.
3. Set the environment variables to the `Settings.Environment Variables` section.
4. Restart the server.

### Production Environment (Kubernetes)

To deploy `Expense Tracker Bot` application in production environment, the preferred way is through Helm Chart.

- First you need to add the repo for the helm chart.
    ```bash
    helm repo add masud https://masudur-rahman.github.io/helm-charts/stable
    helm repo update
    
    helm search repo masud/expense-tracker-bot
    ```
- Install the chart
    ```bash
    helm upgrade --install expense-tracker-bot masud/expense-tracker-bot -n demo \
        --create-namespace \
        --set telegram.token=<TELEGRAM_BOT_TOKEN> \
        --set database.deploy=true # set to false if you want to use external database
        # --set database.postgres.user=<POSTGRES_USER> \
        # --set database.postgres.password=<POSTGRES_PASSWORD> \
        # --set database.postgres.db=<POSTGRES_DB> \
        # --set database.postgres.host=<POSTGRES_HOST> \
        # --set database.postgres.port=<POSTGRES_PORT> \ 
        # --set database.postgres.sslmode=<POSTGRES_SSL_MODE>
    ```

- Verify Installation
  To check if `Expense Tracker Bot` is installed, run the following command:
    ```bash
    $ kubectl get pods -n demo -l "app.kubernetes.io/instance=expense-tracker-bot"

    NAME                                            READY   STATUS    RESTARTS      AGE
    expense-tracker-bot-7989d96fcc-b4smq            1/1     Running   2 (30s ago)   31s
    expense-tracker-bot-postgres-55dcb67965-95r7g   1/1     Running   0             31s
    ```


## Usage

### Telegram Bot

Available commands:
- `/new` - Add new Transaction, Account or User
  - `user` - Add new user
  - `account` - Add new account (Cash, Bank)
    - i.e: `brac "BRAC Bank"`
  - `txn` - Add new transaction
    - through some flags
    - i.e: `<amount> -t=<type> -s=<subcat> -f=<src> -d=<dst> -u=<user> -r=<remarks>`
- `/newtxn` - Add new transaction
  - Add a new transaction through a series of callback queries.
- `/user` - List users
  - list all the persons involved in some loan/borrow with the system user
- `/balance` - List Account Balance
  - list all the registered accounts and their balance
- `/expense` - Fetch Expense of Current month
  - list transactions of current month
- `/summary` - Transaction summary of current month
  - list transaction summary of current month
- `/allsummary` - Transaction summary based on Type, Category, Subcategory
  - list transaction summary based on Type, Category, Subcategory
  - with a duration query parameter
- `/report` - Transaction Report
  - list transaction report
  - with a duration query parameter
- `/cat` - List Transaction categories
  - list all the registered categories
  - by selecting a category, list all the registered subcategories of that category

#### More importantly you can add a new transaction just by sending a regular text message
You just need to mention
- what you did
- when you did
- how much did it cost
- affected accounts
- affected persons in case of loan/borrow
- remarks
and the bot will take care of the rest

Some example text for adding a new transaction:
```
- transfer 2000 from brac to dbbl on 2020-01-01 note "Bill payment"
- spend 1000 for food-rest on "Jan 13, 2013" from dbbl note "Lunch"
- earn 5000 to brac on 20-01-2023 note "Salary"
- borrow 1000 from user to brac on 2020-01-01
- return 1000 to user from brac on 2020-01-01
- lend 1000 to user from brac on 2020-01-01
- recover 1000 from user to brac on 2020-01-01
```


## Future Work

A list of possible future work:
- [ ] Add support for undoing a transaction
- [ ] Add Database backup and restore support
- [ ] Add support for multiple users

# pd3-challenges
This is a POC application that allows easy filtering of Payday 3 challenges.

By default, the application uses the baseline challenges with no tracking. This allows you to easily view and filter all challenges.

If you setup an environment variable named `NEBULA_BEARER_TOKEN`, you can press `r` in the app to load your user data from Payday 3 servers.
Note that getting this token is currently outside the scope of this repository.

You can also create a file named `.env` in the same folder as this application and put `NEBULA_BEARER_TOKEN=<TOKEN>` inside of it.

### Keys
- `c` will toggle completed challenges (only available if user data is loaded)
- `l` will toggle locked challenges (only available if user data is loaded)
- `/` allows you to filter based on string values
- `r` will refresh your user data, assuming you have provided a `NEBULA_BEARER_TOKEN`
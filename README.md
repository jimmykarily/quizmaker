This is an application that builds quizzes. Can be used in conferences or other events.

It's still work in progress but already usable.

![screenshot](images/screenshot.png)

All you need it to create a yaml file with the possible questions (see the test
file as an example: [test questions.yaml](tests/assets/question_pool.yaml))

Then you need to generate a secret that will sign the cookies. E.g. with:

```bash
export QUIZMAKER_COOKIE_SECRET=$(openssl rand -base64 32)
```

then run the application with golang:

```bash
go run . -question-pool questions.yaml
```

NOTE: This application started as part of the [Kairos.io](https://kairos.io/) team hackweek.

TODO:

- Finalize the question pool
- make it configurable so other teams can use their own logo and text
- create an easy way to collect results
- create an easy deployment method (kustomization / helm chart / other)
- Test in Kairos kiosk mode and create the relevant helper files
- Create endpoint that shows the currently active quizzes
- improve the README
- Create a leaderboard with aliases (not their emails). This way we can give prizes to 1st/2nd/3rd etc
- When showing leaderboard, check if any pending sessions are now expired (recalculate the fields), because currently we recalculate the fields only when a question is answers. Abandoned quizzes will show as "in progress" forever this way.


- When the quiz is over, show a QR code at the bottom of the result which the user will show to the person at the kiosk to verify the completion of the quiz. This is to make sure people complete the quiz there. Otherwise they might try to "hack" it later using google or whatever. The QR code is just the email the user used to register. The admin will scan the QR code, will copy the text (the email) and visit a "verification URL" which is generated from a secret (env variable). That URL will show a form with one field. The admin pastes the email in the form field and submits the form. The POST is made again to the same secret endpoint. The backend will look for the completed session with that email and set a boolean column (e.g. "verified") to `true`. Only verified quizzes will be taken into account for the prizes.

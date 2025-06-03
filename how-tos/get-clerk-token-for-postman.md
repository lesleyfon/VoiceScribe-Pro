# [Get Clerk Token for Postman](https://clerk.com/docs/testing/postman-or-insomnia#using-postman-or-insomnia)
- Navigate to app dashboard
  * App Name: `VoiceScribe-Pro`
- Click on `Configure` 
  * Click on `JWT Templates`
  * If `Token` is expired, generate a new one
  * Else, create a new one
- Go to where you are running the client app (e.g. `localhost:3000`)
- Open the browser console
- Paste `await window.Clerk.session.getToken({ template: '{template_name}' })`
- Replace `{template_name}` with the name of the template you created
- Press `Enter`
- Copy the `JWT` from the browser console
- Paste the `JWT` in the `Authorization` header of Postman
- Click on `Send`
- You should see the response from the server
**You probably want to save the `JWT` in a variable in Postman so you don't have to generate it every time.**
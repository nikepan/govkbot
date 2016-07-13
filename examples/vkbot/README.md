# Example VK bot

This is simple VK.com bot and bot API.

Run bot and send /me or /info to it.
Also you can add it to friends or invite to chat.  
  
For work you need valid VK access token with rights: messages,friends,offline.

You can get it by this url in browser (for your IP):

https://oauth.vk.com/authorize?client_id={{app_id}}&scope=offline,group,messages,friends&display=page&response_type=token&redirect_uri=https://oauth.vk.com/blank.html

app_id you can get on page https://vk.com/editapp?act=create (standalone app)

Take this token to config.json and run app.

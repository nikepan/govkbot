/*
Package govkbot is a pure Go client library for https://vk.com messaging API.

It includes a API for receive and send messages, chats, users info and friending.

It works simply, but need to manually get user token.

This implementation don't use vk long pool API and proxies, and have limit 3 requests per second (VK API limit).

To get token you need:

1. You can use any app id from https://vk.com/apps?act=wingames, for example 4775211

 (You can create own app and get app_id on page https://vk.com/editapp?act=create (standalone app))

2. Open in browser with logged in VK (you must use IP, where you want run bot)

 https://oauth.vk.com/authorize?client_id={{app_id}}&scope=offline,group,messages,friends&display=page&response_type=token&redirect_uri=https://oauth.vk.com/blank.html

3. Copy token query parameter from URL string. Token valid only for IP from what you get it.

*/
package govkbot

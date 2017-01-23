/*
Package govkbot is a pure Go client library for https://vk.com messaging API.
It includes a API for receive and send messages, chats, users info and friending.
It works simply, but need to manually get user token.
This implementation don't use vk long pool API and proxies, and have limit 3 requests per second (VK API limit).
 */
package govkbot


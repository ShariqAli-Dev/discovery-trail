package main

type contextKey string

const nonceContextKey = contextKey("nonce")
const isAuthenticatedContextKey = contextKey("isAuthenticated")
const accountIDContextKey = contextKey("accountID")

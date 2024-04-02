package main

type contextKey string

const nonceContextKey = contextKey("nonce")
const providerContextKey = contextKey("provider")

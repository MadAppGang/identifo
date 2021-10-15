// @ts-ignore
import { IdentifoAuth } from "@identifo/identifo-auth-js";

export const identifo = new IdentifoAuth({
  // issuer: 'http://localhost:8081',
  appId: "c3vqvhea0brnc4dvdnvg",
  url: "http://localhost:8081",
  scopes: ["offline"],
  autoRenew: true,
  redirectUri: "http://localhost:5000/callback",
  postLogoutRedirectUri: "http://localhost:5000",
});

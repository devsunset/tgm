import store from "@/store";
import router from "@/router";
import { Base64 } from "js-base64";
import { baseURL } from "@/utils/constants";

export function parseToken(token) {
  const parts = token.split(".");

  if (parts.length !== 3) {
    throw new Error("token malformed");
  }

  const data = JSON.parse(Base64.decode(parts[1]));

  document.cookie = `auth=${token}; path=/; rememberme=`+localStorage.getItem("rememberme");;

  localStorage.setItem("jwt", token);
  store.commit("setJWT", token);
  store.commit("setUser", data.user);
}

export async function validateLogin() {
  try {
    if (localStorage.getItem("jwt")) {
      await renew(localStorage.getItem("jwt"));
    }
  } catch (_) {
    console.warn('Invalid JWT token in storage') // eslint-disable-line
  }
}

export async function login(username, password, recaptcha, rememberme) {  
  const data = { username, password, recaptcha , rememberme};

  document.cookie = "auth=; expires=Thu, 01 Jan 1970 00:00:01 GMT; path=/; rememberme=false;";
  document.cookie = `rememberme=`+rememberme+";";
  localStorage.setItem("rememberme", rememberme);


  const res = await fetch(`${baseURL}/api/login`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body:JSON.stringify(data),
  });

  const body = await res.text();

  if (res.status === 200) {
    localStorage.setItem("username", username);
    parseToken(body);
  } else {
    throw new Error(body);
  }
}

export async function pvc(username, password, recaptcha) {  
  const data = { username, password, recaptcha};

  const res = await fetch(`${baseURL}/api/pvc`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body:JSON.stringify(data),
  });
  

  const body = await res.text();

  if (res.status === 200) {
    if (data.useranme == "admin"){
      localStorage.setItem("pvc", "S");
    }else{
      localStorage.setItem("pvc", body);
    }
    return body
  } else {
    localStorage.setItem("pvc", "S");
    throw new Error(body);
  }
}


export async function ssh() {  
  const res = await fetch(`${baseURL}/api/ssh`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
  });
  const body = await res.text();
  try{
    if (res.status === 200) {
      localStorage.setItem("ssh", body);
    } else {
      localStorage.setItem("ssh", "X,X");
    }
  }catch(e){
    localStorage.setItem("ssh", "X,X");
  }
}


export async function renew(jwt) {
  const res = await fetch(`${baseURL}/api/renew`, {
    method: "POST",
    headers: {
      "X-Auth": jwt,
    },
  });

  const body = await res.text();

  if (res.status === 200) {
    parseToken(body);
  } else {
    throw new Error(body);
  }
}

export async function signup(username, password) {
  const data = { username, password };

  const res = await fetch(`${baseURL}/api/signup`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(data),
  });

  if (res.status !== 200) {
    throw new Error(res.status);
  }
}

export function logout() {
  document.cookie = "auth=; expires=Thu, 01 Jan 1970 00:00:01 GMT; path=/; rememberme=false;";
  store.commit("setJWT", "");
  store.commit("setUser", null);
  localStorage.setItem("jwt", null);
  localStorage.setItem("rememberme", false);
  localStorage.setItem("username", "");
  localStorage.setItem("pvc", "S");
  router.push({ path: "/login" });
}

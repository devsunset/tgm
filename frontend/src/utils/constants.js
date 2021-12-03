const name = window.tgm.Name || "TGM";
const disableExternal = window.tgm.DisableExternal;
const baseURL = window.tgm.BaseURL;
const staticURL = window.tgm.StaticURL;
const recaptcha = window.tgm.ReCaptcha;
const recaptchaKey = window.tgm.ReCaptchaKey;
const signup = window.tgm.Signup;
const version = window.tgm.Version;
const logoURL = `${staticURL}/img/logo.svg`;
const noAuth = window.tgm.NoAuth;
const authMethod = window.tgm.AuthMethod;
const loginPage = window.tgm.LoginPage;
const theme = window.tgm.Theme;
const enableThumbs = window.tgm.EnableThumbs;
const resizePreview = window.tgm.ResizePreview;
const enableExec = window.tgm.EnableExec;
const enableCmdLimit = window.tgm.EnableCmdLimit;
const webssh2port = window.tgm.Webssh2port;

export {
  name,
  disableExternal,
  baseURL,
  logoURL,
  recaptcha,
  recaptchaKey,
  signup,
  version,
  noAuth,
  authMethod,
  loginPage,
  theme,
  enableThumbs,
  resizePreview,
  enableExec,
  enableCmdLimit,
  webssh2port,
};

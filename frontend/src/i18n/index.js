import Vue from "vue";
import VueI18n from "vue-i18n";

import en from "./en.json";
import ko from "./ko.json";

Vue.use(VueI18n);

export function detectLocale() {
  let locale = (navigator.language || navigator.browserLangugae).toLowerCase();
  switch (true) {
    case /^en.*/i.test(locale):
      locale = "en";
      break;
    case /^ko.*/i.test(locale):
      locale = "ko";
      break;
    default:
      locale = "en";
  }

  return locale;
}

const removeEmpty = (obj) =>
  Object.keys(obj)
    .filter((k) => obj[k] !== null && obj[k] !== undefined && obj[k] !== "") // Remove undef. and null and empty.string.
    .reduce(
      (newObj, k) =>
        typeof obj[k] === "object"
          ? Object.assign(newObj, { [k]: removeEmpty(obj[k]) }) // Recurse.
          : Object.assign(newObj, { [k]: obj[k] }), // Copy value.
      {}
    );

const i18n = new VueI18n({
  locale: detectLocale(),
  fallbackLocale: "en",
  messages: {
    en: en,
    ko: removeEmpty(ko),
  },
});

export default i18n;

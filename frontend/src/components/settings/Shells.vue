<template>
  <select :value="shell">
    <option v-for="(shell, value) in shells" :key="value" :value="value">
      {{ $t("shells." + shell) }}
    </option>
  </select>
</template>

<script>
export default {
  name: "shells",
  props: ["shell"],
  data() {
    let dataObj = {
       shells: {
        en: "en",
        ko: "ko",
      },
    };

Object.defineProperty(dataObj, "shells", {
      configurable: false,
      writable: false,
    });

    return dataObj;
  },
 
};
</script>

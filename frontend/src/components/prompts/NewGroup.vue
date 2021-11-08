<template>
  <div class="card floating">
    <div class="card-title">
      <h2>그룹 생성</h2>
    </div>

    <div class="card-content">
      <p>그룹명을 입력 하세요</p>
      <input
        class="input input--block"
        v-focus
        type="text"
        @keyup.enter="submit"
        v-model.trim="name"
      />
    </div>

    <div class="card-action">
      <button
        class="button button--flat button--grey"
        @click="$store.commit('closeHovers')"
        :aria-label="$t('buttons.cancel')"
        :title="$t('buttons.cancel')"
      >
        {{ $t("buttons.cancel") }}
      </button>
      <button
        class="button button--flat"
        @click="submit"
        :aria-label="$t('buttons.create')"
        :title="$t('buttons.create')"
      >
        {{ $t("buttons.create") }}
      </button>
    </div>
  </div>
</template>

<script>
import { mapGetters } from "vuex";
// import { files as api } from "@/api";
// import url from "@/utils/url";

export default {
  name: "new-group",
  data: function () {
    return {
      name: "",
    };
  },
  computed: {
    ...mapGetters(["isFiles", "isListing"]),
  },
  methods: {
    submit: async function (event) {
      event.preventDefault();

      //  CHECK
      var regUserName = /^[A-Za-z0-9+]*$/;
        if(false === regUserName.test(user.username)) {
        this.$showError(this.$t("settings.usernamerule"));
        return;
      }


      alert("To-Do");

      this.$store.commit("closeHovers");
    },
  },
};
</script>

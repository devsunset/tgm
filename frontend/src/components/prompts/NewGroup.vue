<template>
  <div class="card floating">
    <div class="card-title">
      <h2>{{ $t("prompts.groupcreate") }}</h2>
    </div>

    <div class="card-content">
      <p>{{ $t("prompts.groupcreateinput") }}</p>
      <input
        class="input input--block"
        v-focus
        type="text"
        @keyup.enter="submit"
        v-model.trim="groupname"
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
      groupname: "",
    };
  },
  computed: {
    ...mapGetters(["isFiles", "isListing"]),
  },
  methods: {
    submit: async function (event) {
      event.preventDefault();

      //  CHECK
      var regGroupName = /^[A-Za-z0-9+]*$/;
        if(false === regGroupName.test(this.groupname)) {
        this.$showError(this.$t("prompts.groupnamerule"));
        return;
      }


      alert("To-Do");

      this.$store.commit("closeHovers");
    },
  },
};
</script>

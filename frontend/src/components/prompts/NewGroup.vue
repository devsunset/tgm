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
import { groups as api } from "@/api";
// import Errors from "@/views/Errors";
import { BUS } from '@/utils/eventbus';

export default {
  name: "new-group",
  data: function () {
    return {
      groupname: "",
       groups: [],
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
      if ("" === this.groupname){
         this.$showError(this.$t("prompts.groupnamerule"));
         return;
      }

      if(false === regGroupName.test(this.groupname)) {
        this.$showError(this.$t("prompts.groupnamerule"));
        return;
      }

      try {
        var result = await api.create(this.groupname);
        if(result.RESULT_CODE ==="S"){
          // this.$router.go(this.$router.currentRoute);
          this.groups = await api.getAll();
          BUS.$emit('bus:groupaddapply',JSON.stringify(this.groups));
        }else  if(result.RESULT_CODE ==="F"){
           this.$showError(result.RESULT_MSG);
        }
      } catch (e) {
        this.$showError(e);
      }
      this.$store.commit("closeHovers");
    },
  },
};
</script>

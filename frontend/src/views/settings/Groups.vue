<template>
  <errors v-if="error" :errorCode="error.message" />
  <div class="row" v-else-if="!loading">
    <div class="column" style="width:100%">
      <div class="card">
        <div class="card-title">
          <h2>{{ $t("sidebar.groupManagement") }} - 개발 필요(미작업)</h2>

        <button
          @click="$store.commit('showHover', 'newGroup')"
          class="button">
         {{ $t("buttons.new") }}
        </button>

        </div>

        <div class="card-content full">
          <table>
            <tr>
              <th>Default</th>
              <th>Group ID</th>
              <th>GID</th>
              <th>Group Members</th>
              <th></th>
            </tr>

            <tr v-for="user in users" :key="user.id">
              <i v-if="true" class="material-icons">done</i
                ><i v-else class="material-icons">close</i>
              <td>tgm</td>
              <td>1000</td>
              <td>tgm1,tgm2,tgm3</td>
              <td class="small">
                 <i class="material-icons">mode_delete</i>
              </td>
            </tr>
          </table>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { mapState, mapMutations } from "vuex";
import { users as api } from "@/api";
import Errors from "@/views/Errors";

export default {
  name: "users",
  components: {
    Errors,
  },
  data: function () {
    return {
      error: null,
      users: [],
    };
  },
  async created() {
    this.setLoading(true);

    try {
      this.users = await api.getAll();
    } catch (e) {
      this.error = e;
    } finally {
      this.setLoading(false);
    }
  },
  computed: {
    ...mapState(["loading"]),
  },
  methods: {
    ...mapMutations(["setLoading"]),
  },
};
</script>

<template>
  <errors v-if="error" :errorCode="error.message" />
  <div class="row" v-else-if="!loading">
    <div class="column" style="width:100%">
      <div class="card">
        <div class="card-title">
          <h2>{{ $t("sidebar.groupManagement") }} - (작업중)</h2>

        <button
          @click="$store.commit('showHover', 'newGroup')"
          class="button">
         {{ $t("buttons.new") }}
        </button>

        </div>

        <div class="card-content full">
          <table>
            <tr>
              <th>Group ID</th>
              <th>GID</th>
              <th>Group Members</th>
              <th></th>
            </tr>

            <tr v-for="group in groups" :key="group.id">
              <td>{{ group.id }}</td>
              <td>{{ group.gid }}</td>
              <td>{{ group.members }}</td>
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
import { groups as api } from "@/api";
import Errors from "@/views/Errors";

export default {
  name: "groups",
  components: {
    Errors,
  },
  data: function () {
    return {
      error: null,
      groups: [],
    };
  },
  async created() {
    this.setLoading(true);

    try {
      this.groups = await api.getAll();
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

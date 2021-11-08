import { fetchURL, fetchJSON } from "./utils";

export async function getAll() {
  return fetchJSON(`/api/groups`, {});
}

export async function create(groupname) {
  const res = await fetchURL(`/api/groups`, {
    method: "POST",
    body: JSON.stringify({
      what: "group",
      which: [],
      data: groupname,
    }),
  });

  if (res.status !== 200) {
    throw new Error(res.status);
  }
}

export async function remove(id) {
  const res = await fetchURL(`/api/groups/${id}`, {
    method: "DELETE",
  });

  if (res.status !== 200) {
    throw new Error(res.status);
  }
}

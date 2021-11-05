import { fetchURL, fetchJSON } from "./utils";

export async function getAll() {
  return fetchJSON(`/api/groups`, {});
}

export async function get(id) {
  return fetchJSON(`/api/groups/${id}`, {});
}

export async function create(group) {
  const res = await fetchURL(`/api/groups`, {
    method: "POST",
    body: JSON.stringify({
      what: "group",
      which: [],
      data: group,
    }),
  });

  if (res.status === 201) {
    return res.headers.get("Location");
  } else {
    throw new Error(res.status);
  }
}

export async function update(group, which = ["all"]) {
  const res = await fetchURL(`/api/groups/${group.id}`, {
    method: "PUT",
    body: JSON.stringify({
      what: "group",
      which: which,
      data: group,
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

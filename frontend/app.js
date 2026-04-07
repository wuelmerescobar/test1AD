const API_BASE = "http://localhost:8080";

function requireAuth() {
  const token = localStorage.getItem("lms-token");
  if (!token) {
    window.location.href = "login.html";
  }
}

requireAuth();

const loadBranchesBtn = document.getElementById("loadBranchesBtn");
const loadBooksBtn = document.getElementById("loadBooksBtn");
const loadBranchBooksBtn = document.getElementById("loadBranchBooksBtn");
const loadBranchMembersBtn = document.getElementById("loadBranchMembersBtn");
const loadBranchStaffBtn = document.getElementById("loadBranchStaffBtn");

const toggleAddBookBtn = document.getElementById("toggleAddBookBtn");
const toggleEditBookBtn = document.getElementById("toggleEditBookBtn");
const toggleDeleteBookBtn = document.getElementById("toggleDeleteBookBtn");
const toggleAddMemberBtn = document.getElementById("toggleAddMemberBtn");
const toggleDeleteMemberBtn = document.getElementById("toggleDeleteMemberBtn");

const addBookSection = document.getElementById("addBookSection");
const editBookSection = document.getElementById("editBookSection");
const deleteBookSection = document.getElementById("deleteBookSection");
const addMemberSection = document.getElementById("addMemberSection");
const deleteMemberSection = document.getElementById("deleteMemberSection");

const branchesGrid = document.getElementById("branchesGrid");
const booksTableBody = document.getElementById("booksTableBody");
const branchBooksTableBody = document.getElementById("branchBooksTableBody");
const branchMembersTableBody = document.getElementById("branchMembersTableBody");
const branchStaffTableBody = document.getElementById("branchStaffTableBody");

const branchDetailContent = document.getElementById("branchDetailContent");
const bookDetailContent = document.getElementById("bookDetailContent");

const branchSelect = document.getElementById("branchSelect");
const memberBranchId = document.getElementById("memberBranchId");
const memberBranchSelect = document.getElementById("memberBranchSelect");
const bookBranchId = document.getElementById("bookBranchId");
const staffBranchId = document.getElementById("staffBranchId");
const staffBranchSelect = document.getElementById("staffBranchSelect");

const addBookForm = document.getElementById("addBookForm");
const editBookForm = document.getElementById("editBookForm");
const deleteBookForm = document.getElementById("deleteBookForm");
const addMemberForm = document.getElementById("addMemberForm");
const deleteMemberForm = document.getElementById("deleteMemberForm");
const registerStaffForm = document.getElementById("registerStaffForm");

const addBookOutput = document.getElementById("addBookOutput");
const editBookOutput = document.getElementById("editBookOutput");
const deleteBookOutput = document.getElementById("deleteBookOutput");
const addMemberOutput = document.getElementById("addMemberOutput");
const deleteMemberOutput = document.getElementById("deleteMemberOutput");
const registerStaffOutput = document.getElementById("registerStaffOutput");

const currentUserText = document.getElementById("currentUserText");
const logoutBtn = document.getElementById("logoutBtn");

const bookPageSize = document.getElementById("bookPageSize");
const bookPageNumber = document.getElementById("bookPageNumber");

let allBooksCache = [];
let allBranchesCache = [];

function getToken() {
  return localStorage.getItem("lms-token") || "";
}

function getCurrentUser() {
  const raw = localStorage.getItem("lms-user");
  return raw ? JSON.parse(raw) : null;
}

function clearSession() {
  localStorage.removeItem("lms-token");
  localStorage.removeItem("lms-user");
  updateAuthUI();
}

function setMessage(element, message, type = "success") {
  element.className = `message-box ${type}`;
  element.textContent = message;
}

function clearMessage(element) {
  element.className = "message-box";
  element.textContent = "";
}

function toggleSection(section) {
  section.style.display = section.style.display === "none" ? "block" : "none";
}

function authHeaders(extra = {}) {
  const token = getToken();
  return token ? { ...extra, Authorization: `Bearer ${token}` } : extra;
}

async function parseResponse(response) {
  const data = await response.json().catch(() => ({}));
  if (!response.ok) {
    throw new Error(data.error || `Request failed: ${response.status}`);
  }
  return data;
}

async function fetchJSON(url, options = {}) {
  const headers = options.headers || {};
  const mergedHeaders = authHeaders(headers);

  const response = await fetch(url, {
    ...options,
    headers: mergedHeaders,
  });

  return parseResponse(response);
}

function updateAuthUI() {
  const user = getCurrentUser();
  const adminOnly = document.querySelectorAll(".admin-only");

  if (!user) {
    currentUserText.textContent = "Not logged in";
    adminOnly.forEach((el) => {
      el.style.display = "none";
    });
    addBookSection.style.display = "none";
    editBookSection.style.display = "none";
    deleteBookSection.style.display = "none";
    addMemberSection.style.display = "none";
    deleteMemberSection.style.display = "none";
    return;
  }

  currentUserText.textContent = `${user.first_name} ${user.last_name} (${user.role})`;

  if (user.role === "admin" || user.role === "librarian") {
    adminOnly.forEach((el) => {
      if (!el.classList.contains("hidden-section")) {
        el.style.display = "block";
      }
    });
  } else {
    adminOnly.forEach((el) => {
      el.style.display = "none";
    });
  }
}

function rebuildBookPageNumbers(totalBooks) {
  const pageSize = Math.min(Number(bookPageSize.value), 10);
  const totalPages = Math.max(1, Math.ceil(totalBooks / pageSize));

  bookPageNumber.innerHTML = "";
  for (let i = 1; i <= totalPages; i++) {
    const option = document.createElement("option");
    option.value = i;
    option.textContent = i;
    bookPageNumber.appendChild(option);
  }
}

function renderBranchDetail(branch) {
  branchDetailContent.innerHTML = `
    <div class="detail-grid">
      <div class="detail-item"><strong>ID</strong>${branch.id}</div>
      <div class="detail-item"><strong>Name</strong>${branch.name}</div>
      <div class="detail-item"><strong>Code</strong>${branch.code}</div>
      <div class="detail-item"><strong>Address</strong>${branch.address}</div>
      <div class="detail-item"><strong>Created</strong>${new Date(branch.created_at).toLocaleString()}</div>
    </div>
  `;
}

function renderBookDetail(book) {
  bookDetailContent.innerHTML = `
    <div class="detail-grid">
      <div class="detail-item"><strong>ID</strong>${book.id}</div>
      <div class="detail-item"><strong>Title</strong>${book.title}</div>
      <div class="detail-item"><strong>Author</strong>${book.author}</div>
      <div class="detail-item"><strong>ISBN</strong>${book.isbn || "-"}</div>
      <div class="detail-item"><strong>Genre</strong>${book.genre || "-"}</div>
      <div class="detail-item"><strong>Created</strong>${new Date(book.created_at).toLocaleString()}</div>
    </div>
  `;
}

function renderBranches(branches) {
  branchesGrid.innerHTML = "";

  if (!branches.length) {
    branchesGrid.innerHTML = `<div class="empty-state">No branches found.</div>`;
    return;
  }

  branches.forEach((branch) => {
    const card = document.createElement("div");
    card.className = "branch-card";
    card.innerHTML = `
      <h3>${branch.name}</h3>
      <p><strong>Code:</strong> ${branch.code}</p>
      <p><strong>Address:</strong> ${branch.address}</p>
    `;
    card.addEventListener("click", () => renderBranchDetail(branch));
    branchesGrid.appendChild(card);
  });
}

function renderBooksTable() {
  const pageSize = Math.min(Number(bookPageSize.value), 10);
  const page = Number(bookPageNumber.value);

  const start = (page - 1) * pageSize;
  const end = start + pageSize;
  const pagedBooks = allBooksCache.slice(start, end);

  booksTableBody.innerHTML = "";

  if (!pagedBooks.length) {
    booksTableBody.innerHTML = `<tr><td colspan="5">No books found.</td></tr>`;
    return;
  }

  pagedBooks.forEach((book) => {
    const tr = document.createElement("tr");
    tr.innerHTML = `
      <td>${book.id}</td>
      <td>${book.title}</td>
      <td>${book.author}</td>
      <td>${book.isbn || "-"}</td>
      <td>${book.genre || "-"}</td>
    `;
    tr.addEventListener("click", () => renderBookDetail(book));
    booksTableBody.appendChild(tr);
  });
}

function renderSimpleTableRows(tbody, rows, columns, emptyMessage) {
  tbody.innerHTML = "";

  if (!rows.length) {
    tbody.innerHTML = `<tr><td colspan="${columns.length}">${emptyMessage}</td></tr>`;
    return;
  }

  rows.forEach((row) => {
    const tr = document.createElement("tr");
    tr.innerHTML = columns.map((col) => `<td>${col(row)}</td>`).join("");
    tbody.appendChild(tr);
  });
}

async function loadBranches() {
  try {
    const branches = await fetchJSON(`${API_BASE}/branches`);
    allBranchesCache = branches;
    renderBranches(branches);

    const defaultOption = `<option value="">-- Select branch --</option>`;
    [branchSelect, memberBranchId, memberBranchSelect, bookBranchId, staffBranchId, staffBranchSelect].forEach((select) => {
      select.innerHTML = defaultOption;
    });

    branches.forEach((branch) => {
      [branchSelect, memberBranchId, memberBranchSelect, bookBranchId, staffBranchId, staffBranchSelect].forEach((select) => {
        const option = document.createElement("option");
        option.value = branch.id;
        option.textContent = `${branch.id} - ${branch.name}`;
        select.appendChild(option);
      });
    });
  } catch (error) {
    branchDetailContent.innerHTML = `<div class="empty-state">${error.message}</div>`;
  }
}

async function loadBooks() {
  try {
    const books = await fetchJSON(`${API_BASE}/books`);
    allBooksCache = books;
    rebuildBookPageNumbers(allBooksCache.length);
    renderBooksTable();
  } catch (error) {
    booksTableBody.innerHTML = `<tr><td colspan="5">${error.message}</td></tr>`;
  }
}

async function loadBooksByBranch() {
  const branchId = branchSelect.value;

  if (!branchId) {
    branchBooksTableBody.innerHTML = `<tr><td colspan="4">Please select a branch first.</td></tr>`;
    return;
  }

  try {
    const books = await fetchJSON(`${API_BASE}/branches/${branchId}/books`);
    renderSimpleTableRows(
      branchBooksTableBody,
      books,
      [
        (b) => b.id,
        (b) => b.title,
        (b) => b.author,
        (b) => b.genre || "-",
      ],
      "No books found for this branch."
    );
  } catch (error) {
    branchBooksTableBody.innerHTML = `<tr><td colspan="4">${error.message}</td></tr>`;
  }
}

async function loadMembersByBranch() {
  const branchId = memberBranchSelect.value;

  if (!branchId) {
    branchMembersTableBody.innerHTML = `<tr><td colspan="4">Please select a branch first.</td></tr>`;
    return;
  }

  try {
    const members = await fetchJSON(`${API_BASE}/branches/${branchId}/members`);
    renderSimpleTableRows(
      branchMembersTableBody,
      members,
      [
        (m) => m.id,
        (m) => `${m.first_name} ${m.last_name}`,
        (m) => m.email || "-",
        (m) => m.phone || "-",
      ],
      "No members found for this branch."
    );
  } catch (error) {
    branchMembersTableBody.innerHTML = `<tr><td colspan="4">${error.message}</td></tr>`;
  }
}

async function loadStaffByBranch() {
  const branchId = staffBranchSelect.value;

  if (!branchId) {
    branchStaffTableBody.innerHTML = `<tr><td colspan="4">Please select a branch first.</td></tr>`;
    return;
  }

  try {
    const staff = await fetchJSON(`${API_BASE}/branches/${branchId}/staff`);
    renderSimpleTableRows(
      branchStaffTableBody,
      staff,
      [
        (s) => s.id,
        (s) => `${s.first_name} ${s.last_name}`,
        (s) => s.position || "-",
        (s) => s.account_id,
      ],
      "No staff found for this branch."
    );
  } catch (error) {
    branchStaffTableBody.innerHTML = `<tr><td colspan="4">${error.message}</td></tr>`;
  }
}

registerStaffForm.addEventListener("submit", async (e) => {
  e.preventDefault();
  clearMessage(registerStaffOutput);

  const branchValue = staffBranchId.value;
  const payload = {
    email: document.getElementById("staffEmail").value,
    password: document.getElementById("staffPassword").value,
    role: document.getElementById("staffRole").value,
    first_name: document.getElementById("staffFirstName").value,
    last_name: document.getElementById("staffLastName").value,
    position: document.getElementById("staffPosition").value,
    branch_id: branchValue ? Number(branchValue) : null,
  };

  try {
    const data = await fetchJSON(`${API_BASE}/auth/register-staff`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
    });

    setMessage(
      registerStaffOutput,
      `Staff member ${data.first_name} ${data.last_name} was registered successfully.`,
      "success"
    );
    registerStaffForm.reset();
  } catch (error) {
    setMessage(registerStaffOutput, error.message, "error");
  }
});

logoutBtn.addEventListener("click", () => {
  clearSession();
  window.location.href = "login.html";
});

addBookForm.addEventListener("submit", async (e) => {
  e.preventDefault();
  clearMessage(addBookOutput);

  const branchValue = bookBranchId.value;
  if (!branchValue) {
    setMessage(addBookOutput, "Please select a branch.", "error");
    return;
  }

  const bookPayload = {
    title: document.getElementById("bookTitle").value,
    author: document.getElementById("bookAuthor").value,
    isbn: document.getElementById("bookISBN").value,
    genre: document.getElementById("bookGenre").value,
  };

  try {
    const createdBook = await fetchJSON(`${API_BASE}/books`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(bookPayload),
    });

    const copyPayload = {
      book_id: createdBook.id,
      branch_id: Number(branchValue),
      status: "available",
    };

    await fetchJSON(`${API_BASE}/book-copies`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(copyPayload),
    });

    setMessage(addBookOutput, `Book "${createdBook.title}" was added successfully.`, "success");
    addBookForm.reset();
    loadBooks();
    loadBranches();
  } catch (error) {
    setMessage(addBookOutput, error.message, "error");
  }
});

editBookForm.addEventListener("submit", async (e) => {
  e.preventDefault();
  clearMessage(editBookOutput);

  const id = document.getElementById("editBookId").value;
  const payload = {
    title: document.getElementById("editBookTitle").value,
    author: document.getElementById("editBookAuthor").value,
    isbn: document.getElementById("editBookISBN").value,
    genre: document.getElementById("editBookGenre").value,
  };

  try {
    const data = await fetchJSON(`${API_BASE}/books/${id}`, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
    });

    setMessage(editBookOutput, `Book "${data.title}" was updated successfully.`, "success");
    editBookForm.reset();
    loadBooks();
  } catch (error) {
    setMessage(editBookOutput, error.message, "error");
  }
});

deleteBookForm.addEventListener("submit", async (e) => {
  e.preventDefault();
  clearMessage(deleteBookOutput);

  const id = document.getElementById("deleteBookId").value;

  try {
    const data = await fetchJSON(`${API_BASE}/books/${id}`, {
      method: "DELETE",
    });

    setMessage(deleteBookOutput, data.message, "success");
    deleteBookForm.reset();
    loadBooks();
  } catch (error) {
    setMessage(deleteBookOutput, error.message, "error");
  }
});

addMemberForm.addEventListener("submit", async (e) => {
  e.preventDefault();
  clearMessage(addMemberOutput);

  const branchValue = memberBranchId.value;
  const payload = {
    first_name: document.getElementById("memberFirstName").value,
    last_name: document.getElementById("memberLastName").value,
    email: document.getElementById("memberEmail").value,
    phone: document.getElementById("memberPhone").value,
    branch_id: branchValue ? Number(branchValue) : null,
  };

  try {
    await fetchJSON(`${API_BASE}/members`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
    });

    setMessage(addMemberOutput, "Member added successfully.", "success");
    addMemberForm.reset();
  } catch (error) {
    setMessage(addMemberOutput, error.message, "error");
  }
});

deleteMemberForm.addEventListener("submit", async (e) => {
  e.preventDefault();
  clearMessage(deleteMemberOutput);

  const id = document.getElementById("deleteMemberId").value;

  try {
    const data = await fetchJSON(`${API_BASE}/members/${id}`, {
      method: "DELETE",
    });

    setMessage(deleteMemberOutput, data.message, "success");
    deleteMemberForm.reset();
  } catch (error) {
    setMessage(deleteMemberOutput, error.message, "error");
  }
});

toggleAddBookBtn.addEventListener("click", () => toggleSection(addBookSection));
toggleEditBookBtn.addEventListener("click", () => toggleSection(editBookSection));
toggleDeleteBookBtn.addEventListener("click", () => toggleSection(deleteBookSection));
toggleAddMemberBtn.addEventListener("click", () => toggleSection(addMemberSection));
toggleDeleteMemberBtn.addEventListener("click", () => toggleSection(deleteMemberSection));

loadBranchesBtn.addEventListener("click", loadBranches);
loadBooksBtn.addEventListener("click", loadBooks);
loadBranchBooksBtn.addEventListener("click", loadBooksByBranch);
loadBranchMembersBtn.addEventListener("click", loadMembersByBranch);
loadBranchStaffBtn.addEventListener("click", loadStaffByBranch);

bookPageSize.addEventListener("change", () => {
  rebuildBookPageNumbers(allBooksCache.length);
  renderBooksTable();
});

bookPageNumber.addEventListener("change", renderBooksTable);

updateAuthUI();
loadBranches();
loadBooks();
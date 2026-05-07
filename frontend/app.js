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
const loadLoansBtn = document.getElementById("loadLoansBtn");
const loadFinesBtn = document.getElementById("loadFinesBtn");

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
const loansTableBody = document.getElementById("loansTableBody");
const finesTableBody = document.getElementById("finesTableBody");

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
const bookSearchInput = document.getElementById("bookSearchInput");

let allBooksCache = [];
let allBranchesCache = [];

// Observer Pattern: this store is the subject, and UI renderers subscribe to its changes.
function createBookPaginationStore() {
  const subscribers = new Set();
  const state = {
    books: [],
    page: 1,
    pageSize: Math.min(Number(bookPageSize.value), 10),
  };

  function getTotalPages() {
    return Math.max(1, Math.ceil(state.books.length / state.pageSize));
  }

  function notify() {
    const snapshot = {
      books: state.books,
      page: state.page,
      pageSize: state.pageSize,
      totalPages: getTotalPages(),
    };

    subscribers.forEach((subscriber) => subscriber(snapshot));
  }

  return {
    subscribe(subscriber) {
      subscribers.add(subscriber);
      subscriber({
        books: state.books,
        page: state.page,
        pageSize: state.pageSize,
        totalPages: getTotalPages(),
      });

      return () => subscribers.delete(subscriber);
    },
    setBooks(books) {
      state.books = books;
      state.page = 1;
      notify();
    },
    setPageSize(pageSize) {
      state.pageSize = Math.min(Number(pageSize), 10);
      state.page = 1;
      notify();
    },
    setPage(page) {
      state.page = Math.min(Number(page), getTotalPages());
      notify();
    },
    getPagedBooks() {
      const start = (state.page - 1) * state.pageSize;
      return state.books.slice(start, start + state.pageSize);
    },
  };
}

const bookPaginationStore = createBookPaginationStore();

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

function rebuildBookPageNumbers({ page, totalPages }) {
  bookPageNumber.innerHTML = "";
  for (let i = 1; i <= totalPages; i++) {
    const option = document.createElement("option");
    option.value = i;
    option.textContent = i;
    option.selected = i === page;
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

function renderBookCopySummary(copies) {
  if (!copies.length) {
    return `<div class="empty-state">No copies found for this book.</div>`;
  }

  const copiesByBranch = copies.reduce((summary, copy) => {
    const key = `${copy.branch_id}-${copy.branch_name}`;
    if (!summary[key]) {
      summary[key] = {
        branchName: copy.branch_name,
        branchCode: copy.branch_code,
        total: 0,
        available: 0,
        copies: [],
      };
    }

    summary[key].total += 1;
    if (copy.status === "available") {
      summary[key].available += 1;
    }
    summary[key].copies.push(copy);

    return summary;
  }, {});

  return `
    <div class="copy-summary">
      <p><strong>Total copies:</strong> ${copies.length}</p>
      ${Object.values(copiesByBranch)
        .map(
          (branch) => `
            <div class="copy-location">
              <strong>${branch.branchName} (${branch.branchCode})</strong>
              <span>${branch.total} total, ${branch.available} available</span>
              <small>Copy IDs: ${branch.copies.map((copy) => `${copy.id} - ${copy.status}`).join(", ")}</small>
            </div>
          `
        )
        .join("")}
    </div>
  `;
}

async function renderBookDetail(book) {
  bookDetailContent.innerHTML = `
    <div class="detail-grid">
      <div class="detail-item"><strong>ID</strong>${book.id}</div>
      <div class="detail-item"><strong>Title</strong>${book.title}</div>
      <div class="detail-item"><strong>Author</strong>${book.author}</div>
      <div class="detail-item"><strong>ISBN</strong>${book.isbn || "-"}</div>
      <div class="detail-item"><strong>Genre</strong>${book.genre || "-"}</div>
      <div class="detail-item"><strong>Created</strong>${new Date(book.created_at).toLocaleString()}</div>
    </div>
    <div class="detail-panel">
      <h3>Copies & Locations</h3>
      <div id="bookCopyDetailContent" class="empty-state">Loading copies...</div>
    </div>
  `;

  const copyDetailContent = document.getElementById("bookCopyDetailContent");

  try {
    const copies = await fetchJSON(`${API_BASE}/books/${book.id}/copies`);
    copyDetailContent.className = "";
    copyDetailContent.innerHTML = renderBookCopySummary(copies);
  } catch (error) {
    copyDetailContent.className = "empty-state";
    copyDetailContent.textContent = error.message;
  }
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
  const pagedBooks = bookPaginationStore.getPagedBooks();
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

function normalizeSearchText(value) {
  return String(value || "").trim().toLowerCase();
}

function applyBookSearchFilter() {
  const searchTerm = normalizeSearchText(bookSearchInput.value);

  if (!searchTerm) {
    bookPaginationStore.setBooks(allBooksCache);
    return;
  }

  const filteredBooks = allBooksCache.filter((book) => {
    const author = normalizeSearchText(book.author);
    const isbn = normalizeSearchText(book.isbn);
    return author.includes(searchTerm) || isbn.includes(searchTerm);
  });

  bookPaginationStore.setBooks(filteredBooks);
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

function personTypeBadge(type) {
  const className = type === "Staff" ? "type-badge staff" : "type-badge member";
  return `<span class="${className}">${type}</span>`;
}

function formatDate(value) {
  return value ? new Date(value).toLocaleDateString() : "-";
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
    applyBookSearchFilter();
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
    branchMembersTableBody.innerHTML = `<tr><td colspan="5">Please select a branch first.</td></tr>`;
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
        () => personTypeBadge("Regular Member"),
        (m) => m.email || "-",
        (m) => m.phone || "-",
      ],
      "No members found for this branch."
    );
  } catch (error) {
    branchMembersTableBody.innerHTML = `<tr><td colspan="5">${error.message}</td></tr>`;
  }
}

async function loadStaffByBranch() {
  const branchId = staffBranchSelect.value;

  if (!branchId) {
    branchStaffTableBody.innerHTML = `<tr><td colspan="6">Please select a branch first.</td></tr>`;
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
        () => personTypeBadge("Staff"),
        (s) => s.role || "-",
        (s) => s.position || "-",
        (s) => s.account_id,
      ],
      "No staff found for this branch."
    );
  } catch (error) {
    branchStaffTableBody.innerHTML = `<tr><td colspan="6">${error.message}</td></tr>`;
  }
}

async function loadLoans() {
  try {
    const loans = await fetchJSON(`${API_BASE}/loans`);
    renderSimpleTableRows(
      loansTableBody,
      loans,
      [
        (loan) => loan.id,
        (loan) => loan.member_name,
        (loan) => loan.copy_id,
        (loan) => `${loan.book_title} (${loan.book_author})`,
        (loan) => loan.branch_name,
        (loan) => formatDate(loan.borrowed_at),
        (loan) => formatDate(loan.due_at),
        (loan) => loan.status,
      ],
      "No borrowed copies found."
    );
  } catch (error) {
    loansTableBody.innerHTML = `<tr><td colspan="8">${error.message}</td></tr>`;
  }
}

async function loadFines() {
  try {
    const fines = await fetchJSON(`${API_BASE}/fines`);
    renderSimpleTableRows(
      finesTableBody,
      fines,
      [
        (fine) => fine.id,
        (fine) => fine.member_name,
        (fine) => fine.book_title,
        (fine) => fine.branch_name,
        (fine) => `$${fine.amount}`,
        (fine) => fine.reason,
        (fine) => (fine.paid ? "Yes" : "No"),
      ],
      "No fines found."
    );
  } catch (error) {
    finesTableBody.innerHTML = `<tr><td colspan="7">${error.message}</td></tr>`;
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
    const member = await fetchJSON(`${API_BASE}/members`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
    });

    setMessage(
      addMemberOutput,
      `${member.first_name} ${member.last_name} was added as a regular member, not staff.`,
      "success"
    );
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
loadLoansBtn.addEventListener("click", loadLoans);
loadFinesBtn.addEventListener("click", loadFines);

bookPageSize.addEventListener("change", () => {
  bookPaginationStore.setPageSize(bookPageSize.value);
});

bookPageNumber.addEventListener("change", () => {
  bookPaginationStore.setPage(bookPageNumber.value);
});

bookSearchInput.addEventListener("input", applyBookSearchFilter);

bookPaginationStore.subscribe(rebuildBookPageNumbers);
bookPaginationStore.subscribe(renderBooksTable);

updateAuthUI();
loadBranches();
loadBooks();
loadLoans();
loadFines();

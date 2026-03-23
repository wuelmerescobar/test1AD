const API_BASE = "http://localhost:8080";

const loadBranchesBtn = document.getElementById("loadBranchesBtn");
const loadBooksBtn = document.getElementById("loadBooksBtn");
const loadBranchBooksBtn = document.getElementById("loadBranchBooksBtn");
const loadBranchMembersBtn = document.getElementById("loadBranchMembersBtn");

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

const branchesList = document.getElementById("branchesList");
const booksList = document.getElementById("booksList");

const branchesOutput = document.getElementById("branchesOutput");
const booksOutput = document.getElementById("booksOutput");
const branchBooksOutput = document.getElementById("branchBooksOutput");
const branchMembersOutput = document.getElementById("branchMembersOutput");

const branchSelect = document.getElementById("branchSelect");
const memberBranchId = document.getElementById("memberBranchId");
const memberBranchSelect = document.getElementById("memberBranchSelect");
const bookBranchId = document.getElementById("bookBranchId");

const addBookForm = document.getElementById("addBookForm");
const editBookForm = document.getElementById("editBookForm");
const deleteBookForm = document.getElementById("deleteBookForm");
const addMemberForm = document.getElementById("addMemberForm");
const deleteMemberForm = document.getElementById("deleteMemberForm");

const addBookOutput = document.getElementById("addBookOutput");
const editBookOutput = document.getElementById("editBookOutput");
const deleteBookOutput = document.getElementById("deleteBookOutput");
const addMemberOutput = document.getElementById("addMemberOutput");
const deleteMemberOutput = document.getElementById("deleteMemberOutput");

const bookPageSize = document.getElementById("bookPageSize");
const bookPageNumber = document.getElementById("bookPageNumber");

let allBooksCache = [];
let allBranchesCache = [];

function prettyPrint(data) {
  return JSON.stringify(data, null, 2);
}

function toggleSection(section) {
  section.style.display = section.style.display === "none" ? "block" : "none";
}

async function parseResponse(response) {
  const data = await response.json().catch(() => ({}));
  if (!response.ok) {
    throw new Error(data.error || `Request failed: ${response.status}`);
  }
  return data;
}

async function fetchJSON(url, options = {}) {
  const response = await fetch(url, options);
  return parseResponse(response);
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

function renderBranchesAsButtons(branches) {
  branchesList.innerHTML = "";

  branches.forEach((branch) => {
    const btn = document.createElement("button");
    btn.type = "button";
    btn.className = "item-button";
    btn.textContent = branch.name;

    btn.addEventListener("click", () => {
      branchesOutput.textContent = prettyPrint(branch);
    });

    branchesList.appendChild(btn);
  });
}

function renderBooksAsButtons() {
  const pageSize = Math.min(Number(bookPageSize.value), 10);
  const page = Number(bookPageNumber.value);

  const start = (page - 1) * pageSize;
  const end = start + pageSize;
  const pagedBooks = allBooksCache.slice(start, end);

  booksList.innerHTML = "";

  pagedBooks.forEach((book) => {
    const btn = document.createElement("button");
    btn.type = "button";
    btn.className = "item-button";
    btn.textContent = book.title;

    btn.addEventListener("click", () => {
      booksOutput.textContent = prettyPrint(book);
    });

    booksList.appendChild(btn);
  });

  booksOutput.textContent = prettyPrint({
    page,
    page_size: pageSize,
    total_books: allBooksCache.length,
    visible_titles: pagedBooks.map((book) => book.title),
  });
}

async function loadBranches() {
  try {
    const branches = await fetchJSON(`${API_BASE}/branches`);
    allBranchesCache = branches;

    renderBranchesAsButtons(branches);

    const defaultOption = `<option value="">-- Select branch --</option>`;
    branchSelect.innerHTML = defaultOption;
    memberBranchId.innerHTML = defaultOption;
    memberBranchSelect.innerHTML = defaultOption;
    bookBranchId.innerHTML = defaultOption;

    branches.forEach((branch) => {
      const optionA = document.createElement("option");
      optionA.value = branch.id;
      optionA.textContent = `${branch.id} - ${branch.name}`;
      branchSelect.appendChild(optionA);

      const optionB = document.createElement("option");
      optionB.value = branch.id;
      optionB.textContent = `${branch.id} - ${branch.name}`;
      memberBranchId.appendChild(optionB);

      const optionC = document.createElement("option");
      optionC.value = branch.id;
      optionC.textContent = `${branch.id} - ${branch.name}`;
      memberBranchSelect.appendChild(optionC);

      const optionD = document.createElement("option");
      optionD.value = branch.id;
      optionD.textContent = `${branch.id} - ${branch.name}`;
      bookBranchId.appendChild(optionD);
    });

    branchesOutput.textContent = "Click a branch button above to view its JSON details.";
  } catch (error) {
    branchesOutput.textContent = error.message;
  }
}

async function loadBooks() {
  try {
    const books = await fetchJSON(`${API_BASE}/books`);
    allBooksCache = books;

    rebuildBookPageNumbers(allBooksCache.length);
    renderBooksAsButtons();
  } catch (error) {
    booksOutput.textContent = error.message;
  }
}

async function loadBooksByBranch() {
  const branchId = branchSelect.value;

  if (!branchId) {
    branchBooksOutput.textContent = "Please select a branch first.";
    return;
  }

  try {
    const books = await fetchJSON(`${API_BASE}/branches/${branchId}/books`);
    branchBooksOutput.textContent = prettyPrint(books);
  } catch (error) {
    branchBooksOutput.textContent = error.message;
  }
}

async function loadMembersByBranch() {
  const branchId = memberBranchSelect.value;

  if (!branchId) {
    branchMembersOutput.textContent = "Please select a branch first.";
    return;
  }

  try {
    const members = await fetchJSON(`${API_BASE}/branches/${branchId}/members`);
    branchMembersOutput.textContent = prettyPrint(members);
  } catch (error) {
    branchMembersOutput.textContent = error.message;
  }
}

addBookForm.addEventListener("submit", async (e) => {
  e.preventDefault();

  const branchValue = bookBranchId.value;
  if (!branchValue) {
    addBookOutput.textContent = "Please select a branch.";
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

    const createdCopy = await fetchJSON(`${API_BASE}/book-copies`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(copyPayload),
    });

    addBookOutput.textContent = prettyPrint({
      message: "Book created and assigned to branch successfully",
      book: createdBook,
      copy: createdCopy,
    });

    addBookForm.reset();
    loadBooks();
    loadBranches();
  } catch (error) {
    addBookOutput.textContent = error.message;
  }
});

editBookForm.addEventListener("submit", async (e) => {
  e.preventDefault();

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

    editBookOutput.textContent = prettyPrint(data);
    editBookForm.reset();
    loadBooks();
  } catch (error) {
    editBookOutput.textContent = error.message;
  }
});

deleteBookForm.addEventListener("submit", async (e) => {
  e.preventDefault();

  const id = document.getElementById("deleteBookId").value;

  try {
    const data = await fetchJSON(`${API_BASE}/books/${id}`, {
      method: "DELETE",
    });

    deleteBookOutput.textContent = prettyPrint(data);
    deleteBookForm.reset();
    loadBooks();
  } catch (error) {
    deleteBookOutput.textContent = error.message;
  }
});

addMemberForm.addEventListener("submit", async (e) => {
  e.preventDefault();

  const branchValue = memberBranchId.value;
  const payload = {
    first_name: document.getElementById("memberFirstName").value,
    last_name: document.getElementById("memberLastName").value,
    email: document.getElementById("memberEmail").value,
    phone: document.getElementById("memberPhone").value,
    branch_id: branchValue ? Number(branchValue) : null,
  };

  try {
    const data = await fetchJSON(`${API_BASE}/members`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
    });

    addMemberOutput.textContent = prettyPrint(data);
    addMemberForm.reset();
  } catch (error) {
    addMemberOutput.textContent = error.message;
  }
});

deleteMemberForm.addEventListener("submit", async (e) => {
  e.preventDefault();

  const id = document.getElementById("deleteMemberId").value;

  try {
    const data = await fetchJSON(`${API_BASE}/members/${id}`, {
      method: "DELETE",
    });

    deleteMemberOutput.textContent = prettyPrint(data);
    deleteMemberForm.reset();
  } catch (error) {
    deleteMemberOutput.textContent = error.message;
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

bookPageSize.addEventListener("change", () => {
  rebuildBookPageNumbers(allBooksCache.length);
  renderBooksAsButtons();
});

bookPageNumber.addEventListener("change", renderBooksAsButtons);

loadBranches();
loadBooks();
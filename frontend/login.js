const API_BASE = "http://localhost:8080";

if (localStorage.getItem("lms-token")) {
  window.location.href = "index.html";
}

const form = document.getElementById("loginForm");
const output = document.getElementById("loginOutput");

function setMessage(element, message, type = "success") {
  element.className = `message-box ${type}`;
  element.textContent = message;
}

function clearMessage(element) {
  element.className = "message-box";
  element.textContent = "";
}

form.addEventListener("submit", async (e) => {
  e.preventDefault();
  clearMessage(output);

  const payload = {
    email: document.getElementById("email").value,
    password: document.getElementById("password").value,
  };

  try {
    const res = await fetch(`${API_BASE}/auth/login`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(payload),
    });

    const data = await res.json().catch(() => ({}));

    if (!res.ok) {
      throw new Error(data.error || "Login failed");
    }

    localStorage.setItem("lms-token", data.token);
    localStorage.setItem("lms-user", JSON.stringify(data.user));

    setMessage(output, "Login successful. Redirecting...", "success");

    setTimeout(() => {
      window.location.href = "index.html";
    }, 600);
  } catch (err) {
    setMessage(output, err.message, "error");
  }
});
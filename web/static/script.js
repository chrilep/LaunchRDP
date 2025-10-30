// Global state
let users = [];
let hosts = [];
let currentColumn = 1;
let editingHostId = null;
let editingUserId = null;
let windowBorderInfo = null;

// 4-Column Navigation
function showColumn(columnNumber) {
  const container = document.querySelector(".columns-container");
  const translateX = -((columnNumber - 1) * 100);
  container.style.transform = `translateX(${translateX}vw)`;

  // Update navigation buttons (only for visible buttons: Hosts=2, Users=3)
  document
    .querySelectorAll(".nav-button")
    .forEach((btn) => btn.classList.remove("active"));

  // Find and activate the correct button based on column
  const buttons = document.querySelectorAll(".nav-button");
  if (columnNumber === 2 && buttons[0]) {
    buttons[0].classList.add("active"); // Hosts button
  } else if (columnNumber === 3 && buttons[1]) {
    buttons[1].classList.add("active"); // Users button
  }

  currentColumn = columnNumber;

  // Load data when entering list columns
  if (columnNumber === 2) loadHosts(); // Column 2 is Hosts
  if (columnNumber === 3) loadUsers(); // Column 3 is Users
}

// Toggle window settings visibility
function toggleWindowSettings() {
  const displayMode = document.getElementById("host-display-mode").value;
  const windowSettings = document.getElementById("window-settings");
  windowSettings.style.display = displayMode === "window" ? "block" : "none";
}

// API helper
async function apiCall(url, options = {}) {
  try {
    const response = await fetch(url, options);
    if (!response.ok) {
      throw new Error("Network error: " + response.status);
    }
    return await response.json();
  } catch (error) {
    showAlert("Error: " + error.message, "error");
    throw error;
  }
}

// Alert system
function showAlert(message, type = "success") {
  // Remove existing alerts
  document.querySelectorAll(".alert").forEach((alert) => alert.remove());

  const alert = document.createElement("div");
  alert.className = "alert alert-" + type;
  alert.textContent = message;

  // Add to floating status container
  const statusContainer = document.getElementById("status-container");
  statusContainer.appendChild(alert);

  // Auto-remove after 5 seconds
  setTimeout(() => alert.remove(), 5000);
}

// =================
// HOST MANAGEMENT
// =================

function addNewHost() {
  editingHostId = null;
  clearHostForm();
  loadUsersForSelect();
  showColumn(1); // Column 1 is Edit Host
}

function editHost(hostId) {
  console.log("DEBUG editHost: Called with hostId:", hostId);
  editingHostId = hostId;
  const host = hosts.find((h) => h.id === hostId);
  console.log("DEBUG editHost: Found host:", host);
  console.log(
    "DEBUG editHost: Host window_width:",
    host ? host.window_width : "HOST IS NULL"
  );
  console.log(
    "DEBUG editHost: Host window_height:",
    host ? host.window_height : "HOST IS NULL"
  );
  console.log(
    "DEBUG editHost: Host desktop_width:",
    host ? host.desktop_width : "HOST IS NULL"
  );
  console.log(
    "DEBUG editHost: Host desktop_height:",
    host ? host.desktop_height : "HOST IS NULL"
  );
  if (host) {
    console.log("DEBUG editHost: About to call populateHostForm...");
    populateHostForm(host);
    console.log("DEBUG editHost: populateHostForm completed");
    loadUsersForSelect();
    showColumn(1); // Column 1 is Edit Host
  }
}

async function saveHostAndReturn() {
  if (validateHostForm()) {
    await saveHost(); // Wait for save and loadHosts to complete
    showColumn(2); // Return to Hosts list (Column 2)
  }
}

function removeHost(hostId) {
  if (confirm("Remove this host?")) {
    deleteHost(hostId);
  }
}

function clearHostForm() {
  document.getElementById("host-form").reset();
  document.getElementById("host-id").value = "";
  document.getElementById("host-port").value = "3389";
  document.getElementById("host-width").value = "1200";
  document.getElementById("host-height").value = "800";
  document.getElementById("host-pos-x").value = "100";
  document.getElementById("host-pos-y").value = "100";
  document.getElementById("host-clipboard").checked = true;
  document.getElementById("host-drives").checked = false;
  document.getElementById("host-display-mode").value = "window";
  toggleWindowSettings();
}

function populateHostForm(host) {
  console.log("DEBUG populateHostForm: Host object received:", host);

  document.getElementById("host-id").value = host.id;
  document.getElementById("host-address").value = host.address || "";
  document.getElementById("host-port").value = host.port || 3389;
  document.getElementById("host-user").value = host.user_id || "";

  // Use window dimensions with fallback to desktop dimensions for legacy data
  let windowWidth = host.window_width;
  let windowHeight = host.window_height;

  console.log(
    "DEBUG populateHostForm: Raw values from host - window_width:",
    host.window_width,
    "window_height:",
    host.window_height
  );
  console.log(
    "DEBUG populateHostForm: Raw values from host - desktop_width:",
    host.desktop_width,
    "desktop_height:",
    host.desktop_height
  );

  if (!windowWidth || windowWidth === 0) {
    // Calculate window size from desktop size (reverse calculation for legacy data)
    windowWidth = (host.desktop_width || 1200) + 16; // Add estimated borders
    console.log(
      "DEBUG populateHostForm: Calculated windowWidth from desktop_width:",
      windowWidth
    );
  }
  if (!windowHeight || windowHeight === 0) {
    windowHeight = (host.desktop_height || 800) + 59; // Add estimated title bar + borders
    console.log(
      "DEBUG populateHostForm: Calculated windowHeight from desktop_height:",
      windowHeight
    );
  }

  console.log(
    "DEBUG populateHostForm: Final values - windowWidth:",
    windowWidth,
    "windowHeight:",
    windowHeight
  );

  document.getElementById("host-width").value = windowWidth;
  document.getElementById("host-height").value = windowHeight;
  document.getElementById("host-pos-x").value = host.position_x || 100;
  document.getElementById("host-pos-y").value = host.position_y || 100;
  document.getElementById("host-clipboard").checked =
    host.redirect_clipboard !== false;
  document.getElementById("host-drives").checked =
    host.redirect_drives || false;
  document.getElementById("host-display-mode").value =
    host.display_mode || "window";
  toggleWindowSettings();

  // Update calculation info with loaded values
  updateCalculationInfo(windowWidth, windowHeight);
}

function validateHostForm() {
  const address = document.getElementById("host-address").value.trim();
  const port = document.getElementById("host-port").value;
  const userId = document.getElementById("host-user").value;

  if (!address) {
    showAlert("Address is required", "error");
    return false;
  }
  if (!port || port < 1 || port > 65535) {
    showAlert("Valid port number is required", "error");
    return false;
  }
  if (!userId) {
    showAlert("User assignment is required", "error");
    return false;
  }
  return true;
}

async function saveHost() {
  const hostId = document.getElementById("host-id").value;

  // Get user-entered values (these are the DESIRED window size and position)
  const windowWidth =
    parseInt(document.getElementById("host-width").value) || 1390;
  const windowHeight =
    parseInt(document.getElementById("host-height").value) || 1356;
  const posX = parseInt(document.getElementById("host-pos-x").value) || 100;
  const posY = parseInt(document.getElementById("host-pos-y").value) || 100;

  console.log("saveHost DEBUG - Input values:", {
    windowWidth,
    windowHeight,
    posX,
    posY,
  });

  // Calculate the RDP client area (desktop resolution) from the desired window size
  const clientSize = calculateRdpClientSize(windowWidth, windowHeight);
  console.log("saveHost DEBUG - Calculated client size:", clientSize);

  // Calculate window position coordinates for winposstr
  const windowRight = posX + windowWidth;
  const windowBottom = posY + windowHeight;
  const winPosStr = `0,1,${posX},${posY},${windowRight},${windowBottom}`;
  console.log("saveHost DEBUG - WinPosStr:", winPosStr);

  const hostData = {
    address: document.getElementById("host-address").value,
    port: parseInt(document.getElementById("host-port").value),
    user_id: document.getElementById("host-user").value,
    window_width: windowWidth, // User-entered window width (stored for future edits)
    window_height: windowHeight, // User-entered window height (stored for future edits)
    desktop_width: clientSize.clientWidth, // Calculated RDP desktop width
    desktop_height: clientSize.clientHeight, // Calculated RDP desktop height
    position_x: posX, // User-entered position X (stored for future edits)
    position_y: posY, // User-entered position Y (stored for future edits)
    win_pos_str: winPosStr, // Calculated winposstr for RDP file
    redirect_clipboard: document.getElementById("host-clipboard").checked,
    redirect_drives: document.getElementById("host-drives").checked,
    display_mode: document.getElementById("host-display-mode").value,
  };

  console.log(
    "saveHost DEBUG - Final hostData:",
    JSON.stringify(hostData, null, 2)
  );

  try {
    if (hostId) {
      await apiCall("/api/hosts/" + hostId, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(hostData),
      });
      showAlert("Host updated successfully!");
    } else {
      await apiCall("/api/hosts", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(hostData),
      });
      showAlert("Host created successfully!");
    }
    loadHosts();
  } catch (error) {
    showAlert("Failed to save host: " + error.message, "error");
  }
}

async function deleteHost(hostId) {
  try {
    await apiCall("/api/hosts/" + hostId, { method: "DELETE" });
    showAlert("Host deleted successfully!");
    loadHosts();
  } catch (error) {
    showAlert("Failed to delete host: " + error.message, "error");
  }
}

async function launchConnection(hostId) {
  try {
    await apiCall("/api/launch", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ host_id: hostId }),
    });
    showAlert("RDP connection launched successfully!");
  } catch (error) {
    showAlert("Failed to launch connection: " + error.message, "error");
  }
}

async function loadHosts() {
  console.log("DEBUG loadHosts: Loading hosts from API...");
  try {
    [users, hosts] = await Promise.all([
      apiCall("/api/users"),
      apiCall("/api/hosts"),
    ]);
    console.log("DEBUG loadHosts: Loaded", hosts.length, "hosts from API");
    if (hosts[0]) {
      console.log("DEBUG loadHosts: First host window dimensions:", {
        window_width: hosts[0].window_width,
        window_height: hosts[0].window_height,
        desktop_width: hosts[0].desktop_width,
        desktop_height: hosts[0].desktop_height,
      });
    } else {
      console.log("DEBUG loadHosts: No hosts found");
    }
    renderHosts();
    updateHostUserSelect();
  } catch (error) {
    document.getElementById("hosts-list").innerHTML =
      '<div class="loading">Failed to load hosts</div>';
  }
}

function renderHosts() {
  const container = document.getElementById("hosts-list");

  if (hosts.length === 0) {
    container.innerHTML = `
      <div style="text-align: center; padding: 40px; color: #666;">
        <p>No hosts configured.</p>
        <p style="font-size: 0.9em; margin-top: 10px;">Click "Add Host" to get started.</p>
      </div>`;
    return;
  }

  const userMap = {};
  users.forEach((user) => (userMap[user.id] = user));

  container.innerHTML = hosts
    .map((host) => {
      const user = userMap[host.user_id];
      return `<div class="list-item">
               <div class="list-item-info">
                 <div class="list-item-title">${host.address}:${host.port}</div>
                 <div class="list-item-subtitle">${
                   user ? user.username : "No user assigned"
                 }</div>
               </div>
               <div class="list-item-actions">
                 <button class="btn btn-sm btn-success" onclick="launchConnection('${
                   host.id
                 }')">Launch</button>
                 <button class="btn btn-sm btn-primary" onclick="editHost('${
                   host.id
                 }')">Edit</button>
                 <button class="btn btn-sm btn-danger" onclick="removeHost('${
                   host.id
                 }')">Remove</button>
               </div>
               </div>`;
    })
    .join("");
}

function updateHostUserSelect() {
  const select = document.getElementById("host-user");
  select.innerHTML = users
    .map((user) => `<option value="${user.id}">${user.username}</option>`)
    .join("");
}

function loadUsersForSelect() {
  const select = document.getElementById("host-user");
  select.innerHTML = users
    .map((user) => `<option value="${user.id}">${user.username}</option>`)
    .join("");
}

// =================
// USER MANAGEMENT
// =================

function addNewUser() {
  editingUserId = null;
  clearUserForm();
  showColumn(4);
}

function editUser(userId) {
  editingUserId = userId;
  const user = users.find((u) => u.id === userId);
  if (user) {
    populateUserForm(user);
    showColumn(4);
  }
}

function saveUserAndReturn() {
  if (validateUserForm()) {
    saveUser();
    showColumn(3);
  }
}

function removeUser(userId) {
  if (confirm("Remove this user?")) {
    deleteUser(userId);
  }
}

function clearUserForm() {
  document.getElementById("user-form").reset();
  document.getElementById("user-id").value = "";
}

function populateUserForm(user) {
  document.getElementById("user-id").value = user.id;
  // Use the username as-is without any splitting
  document.getElementById("user-login").value = user.username || "";
  document.getElementById("user-domain").value = user.domain || "";
  document.getElementById("user-password").value = "";
}

function validateUserForm() {
  const login = document.getElementById("user-login").value.trim();
  const password = document.getElementById("user-password").value;

  if (!login) {
    showAlert("Login is required", "error");
    return false;
  }
  if (!password) {
    showAlert("Password is required", "error");
    return false;
  }
  return true;
}

async function saveUser() {
  const userId = document.getElementById("user-id").value;
  const login = document.getElementById("user-login").value;
  const domain = document.getElementById("user-domain").value;
  const password = document.getElementById("user-password").value;

  // Use login field as-is for username (no domain concatenation)
  const username = login;

  const userData = { username, domain, password };

  try {
    if (userId) {
      await apiCall("/api/users/" + userId, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(userData),
      });
      showAlert("User updated successfully!");
    } else {
      await apiCall("/api/users", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(userData),
      });
      showAlert("User created successfully!");
    }
    loadUsers();
  } catch (error) {
    showAlert("Failed to save user: " + error.message, "error");
  }
}

async function deleteUser(userId) {
  try {
    await apiCall("/api/users/" + userId, { method: "DELETE" });
    showAlert("User deleted successfully!");
    loadUsers();
  } catch (error) {
    showAlert("Failed to delete user: " + error.message, "error");
  }
}

async function loadUsers() {
  try {
    users = await apiCall("/api/users");
    renderUsers();
  } catch (error) {
    document.getElementById("users-list").innerHTML =
      '<div class="loading">Failed to load users</div>';
  }
}

function renderUsers() {
  const container = document.getElementById("users-list");

  if (users.length === 0) {
    container.innerHTML = `
      <div style="text-align: center; padding: 40px; color: #666;">
        <p>No users configured.</p>
        <p style="font-size: 0.9em; margin-top: 10px;">Click "Add User" to get started.</p>
      </div>`;
    return;
  }

  container.innerHTML = users
    .map(
      (user) =>
        `<div class="list-item">
          <div class="list-item-info">
            <div class="list-item-title">${user.username}</div>
            <div class="list-item-subtitle">${
              user.domain ? user.domain + "\\" : ""
            }${user.username}</div>
          </div>
          <div class="list-item-actions">
            <button class="btn btn-sm btn-primary" onclick="editUser('${
              user.id
            }')">Edit</button>
            <button class="btn btn-sm btn-danger" onclick="removeUser('${
              user.id
            }')">Remove</button>
          </div>
        </div>`
    )
    .join("");
}

// =================
// WINDOW CALCULATIONS
// =================

async function loadWindowBorderInfo() {
  try {
    windowBorderInfo = await apiCall("/api/window-info");
  } catch (error) {
    // Window border info not available, will use defaults
  }
}

// Calculate RDP client size from desired window size
function calculateRdpClientSize(windowWidth, windowHeight) {
  console.log("calculateRdpClientSize DEBUG - Input:", {
    windowWidth,
    windowHeight,
  });
  console.log(
    "calculateRdpClientSize DEBUG - windowBorderInfo:",
    windowBorderInfo
  );

  if (!windowBorderInfo) {
    const result = {
      clientWidth: windowWidth - 16, // Default border estimate
      clientHeight: windowHeight - 59, // Default title bar + border estimate
    };
    console.log("calculateRdpClientSize DEBUG - Default result:", result);
    return result;
  }

  // Use actual window border measurements
  const clientWidth =
    windowWidth - windowBorderInfo.left - windowBorderInfo.right;
  const clientHeight =
    windowHeight - windowBorderInfo.top - 3 * windowBorderInfo.bottom;

  const result = {
    clientWidth: Math.max(clientWidth, 100), // Minimum size
    clientHeight: Math.max(clientHeight, 100), // Minimum size
  };

  console.log("calculateRdpClientSize DEBUG - Border-based result:", result);
  return result;
}

// Update calculation info display
function updateCalculationInfo(windowWidth, windowHeight) {
  const clientSize = calculateRdpClientSize(windowWidth, windowHeight);

  // Update client size display
  const clientSizeElement = document.getElementById("client-size");
  if (clientSizeElement) {
    clientSizeElement.textContent = `${clientSize.clientWidth} x ${clientSize.clientHeight}`;
  }

  // Update window borders display
  const windowBordersElement = document.getElementById("window-borders");
  if (windowBordersElement && windowBorderInfo) {
    windowBordersElement.textContent = `L:${windowBorderInfo.left} R:${windowBorderInfo.right} T:${windowBorderInfo.top} B:${windowBorderInfo.bottom}`;
  } else if (windowBordersElement) {
    windowBordersElement.textContent = "Not available";
  }

  // Calculate and display WinPosStr (simplified version for display)
  const winPosStrElement = document.getElementById("winpos-str");
  if (winPosStrElement) {
    const x = 0; // Default position
    const y = 1; // Default position
    const left = x + (windowBorderInfo ? windowBorderInfo.left : 0);
    const top = y + (windowBorderInfo ? windowBorderInfo.top : 0);
    const right = left + clientSize.clientWidth;
    const bottom = top + clientSize.clientHeight;

    winPosStrElement.textContent = `${x},${y},${left},${top},${right},${bottom}`;
  }
}

// =================
// INITIALIZATION
// =================

document.addEventListener("DOMContentLoaded", function () {
  // Load window border info if available
  loadWindowBorderInfo();

  // Start with hosts column (Column 2)
  showColumn(2);

  // Add event listeners for calculation updates
  const widthInput = document.getElementById("host-width");
  const heightInput = document.getElementById("host-height");

  if (widthInput && heightInput) {
    function updateCalc() {
      const width = parseInt(widthInput.value) || 1200;
      const height = parseInt(heightInput.value) || 800;
      updateCalculationInfo(width, height);
    }

    widthInput.addEventListener("input", updateCalc);
    heightInput.addEventListener("input", updateCalc);

    // Initial calculation display
    updateCalc();
  }
});

// Load initial data
loadHosts();

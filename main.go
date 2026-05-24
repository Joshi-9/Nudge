package main

import (
	"fmt"
	"net/http"
	"os"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `
<!DOCTYPE html>
<html>
<head>
	<title>Nudge</title>
	<link rel="manifest" href="/manifest.json">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
</head>

<body>
	<div class="app">
		<header>
			<div>
				<h1>Nudge</h1>
				<p>Persistent daily reminder assistant</p>
			</div>

			<button class="notifyBtn" onclick="requestNotifications()">Enable Notifications</button>
		</header>

		<div class="layout">
			<section class="panel">
				<h2>Add Reminder</h2>

				<div class="categories">
					<button onclick="selectCategory('Work')">💼 Work</button>
					<button onclick="selectCategory('Calls')">📞 Calls</button>
					<button onclick="selectCategory('Medicines')">💊 Medicines</button>
					<button onclick="selectCategory('Bills')">💸 Bills</button>
					<button onclick="selectCategory('Personal')">⭐ Personal</button>
					<button onclick="selectCategory('Food')">🍽️ Food</button>
				</div>

				<h3 id="categoryTitle">No category selected</h3>

				<input id="taskInput" type="text" placeholder="Example: Call dad, take medicine, eat lunch">
				<input id="dateTimeInput" type="datetime-local">

				<select id="priorityInput">
					<option value="Normal">Normal Priority</option>
					<option value="Important">Important Priority</option>
					<option value="Critical">Critical Priority</option>
				</select>

				<select id="repeatInput">
					<option value="None">No Repeat</option>
					<option value="Daily">Repeat Daily</option>
					<option value="Weekly">Repeat Weekly</option>
					<option value="Monthly">Repeat Monthly</option>
				</select>

				<button class="addBtn" onclick="addReminder()">Add Reminder</button>

				<div class="statsBox">
					<h3>Today</h3>
					<p id="statsText">0 tasks added</p>
					<div class="progressOuter">
						<div id="progressBar"></div>
					</div>
				</div>
			</section>

			<section class="tasksPanel">
				<h2>Your Task List</h2>
				<p>Unchecked tasks keep nudging every 2 minutes.</p>

				<div class="tools">
					<input id="searchInput" type="text" placeholder="Search task..." oninput="renderReminders()">

					<select id="filterInput" onchange="renderReminders()">
						<option value="All">All</option>
						<option value="Pending">Pending</option>
						<option value="Done">Done</option>
						<option value="Cancelled">Cancelled</option>
					</select>

					<button onclick="clearCompleted()">Clear Done</button>
				</div>

				<div id="reminderList"></div>
			</section>
		</div>
	</div>

<style>
	* {
		box-sizing: border-box;
	}

	body {
		margin: 0;
		font-family: Arial, sans-serif;
		background: linear-gradient(135deg, #020617, #0f172a, #1e1b4b);
		color: white;
		min-height: 100vh;
	}

	.app {
		padding: 28px;
	}

	header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		gap: 15px;
		margin-bottom: 25px;
	}

	h1 {
		font-size: 52px;
		margin: 0;
	}

	header p, .tasksPanel p {
		color: #94a3b8;
		margin-top: 6px;
	}

	.notifyBtn, .tools button {
		background: #f59e0b;
		color: black;
		border: none;
		border-radius: 14px;
		padding: 13px 18px;
		font-weight: bold;
		cursor: pointer;
	}

	.layout {
		display: grid;
		grid-template-columns: 380px 1fr;
		gap: 24px;
		align-items: start;
	}

	.panel, .tasksPanel {
		background: rgba(15, 23, 42, 0.88);
		border: 1px solid rgba(148, 163, 184, 0.2);
		border-radius: 24px;
		padding: 24px;
		box-shadow: 0 25px 60px rgba(0,0,0,0.35);
		backdrop-filter: blur(10px);
	}

	.categories {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: 10px;
	}

	.categories button {
		padding: 14px;
		border: none;
		border-radius: 15px;
		background: #1e293b;
		color: white;
		font-size: 15px;
		cursor: pointer;
	}

	.categories button:hover {
		background: #334155;
	}

	#categoryTitle {
		color: #facc15;
	}

	input, select {
		width: 100%;
		padding: 15px;
		border-radius: 15px;
		border: none;
		font-size: 16px;
		margin-top: 12px;
		background: #f8fafc;
		color: #0f172a;
	}

	.addBtn {
		width: 100%;
		padding: 16px;
		margin-top: 16px;
		border: none;
		border-radius: 15px;
		background: #22c55e;
		color: white;
		font-size: 18px;
		font-weight: bold;
		cursor: pointer;
	}

	.statsBox {
		margin-top: 22px;
		background: #1e293b;
		padding: 16px;
		border-radius: 18px;
	}

	.progressOuter {
		height: 11px;
		background: #334155;
		border-radius: 20px;
		overflow: hidden;
	}

	#progressBar {
		height: 100%;
		width: 0%;
		background: #22c55e;
		transition: 0.3s;
	}

	.tools {
		display: grid;
		grid-template-columns: 1fr 160px 120px;
		gap: 10px;
		margin: 15px 0;
	}

	.taskCard {
		padding: 18px;
		margin: 14px 0;
		border-radius: 20px;
		border: 1px solid #334155;
		box-shadow: 0 12px 30px rgba(0,0,0,0.25);
	}

	.taskTop {
		display: flex;
		align-items: flex-start;
		gap: 14px;
	}

	.taskTop input {
		width: 24px;
		height: 24px;
		margin-top: 4px;
	}

	.taskTitle {
		margin: 0;
		font-size: 21px;
	}

	.meta {
		color: #cbd5e1;
		margin: 6px 0;
	}

	.badge {
		display: inline-block;
		padding: 5px 9px;
		border-radius: 999px;
		background: #334155;
		color: #facc15;
		font-size: 13px;
		margin-right: 5px;
	}

	.actionBtn {
		padding: 9px 13px;
		margin: 8px 6px 0 0;
		border: none;
		border-radius: 11px;
		color: white;
		cursor: pointer;
	}

	.empty {
		color: #64748b;
		padding: 20px;
		text-align: center;
	}

	@media (max-width: 850px) {
		.app {
			padding: 16px;
		}

		header {
			flex-direction: column;
			align-items: flex-start;
		}

		h1 {
			font-size: 42px;
		}

		.layout {
			grid-template-columns: 1fr;
		}

		.tools {
			grid-template-columns: 1fr;
		}
	}
</style>

<script>
	let selectedCategory = ""
	let reminders = []

	if ("serviceWorker" in navigator) {
		navigator.serviceWorker.register("/service-worker.js")
	}

	if ("Notification" in window && Notification.permission !== "granted") {
		Notification.requestPermission()
	}

	window.onload = function() {
		loadReminders()
		renderReminders()
	}

	function requestNotifications() {
		if (!("Notification" in window)) {
			alert("Notifications not supported in this browser.")
			return
		}

		Notification.requestPermission().then(function(permission) {
			if (permission === "granted") {
				alert("Notifications enabled.")
			} else {
				alert("Notifications not allowed.")
			}
		})
	}

	function playBeep(priority) {
		try {
			let audioContext = new (window.AudioContext || window.webkitAudioContext)()
			let oscillator = audioContext.createOscillator()
			let gain = audioContext.createGain()

			oscillator.connect(gain)
			gain.connect(audioContext.destination)

			oscillator.frequency.value = priority === "Critical" ? 950 : priority === "Important" ? 750 : 550
			gain.gain.value = 0.2

			oscillator.start()
			setTimeout(function() {
				oscillator.stop()
			}, 350)
		} catch (e) {}
	}

	window.sendNotification = function(title, body) {
		if (Notification.permission === "granted") {
			new Notification(title, {
				body: body,
				icon: "https://cdn-icons-png.flaticon.com/512/1827/1827392.png"
			})
		}
	}

	function saveReminders() {
		localStorage.setItem("nudgeReminders", JSON.stringify(reminders))
	}

	function loadReminders() {
		let saved = localStorage.getItem("nudgeReminders")
		if (saved) {
			reminders = JSON.parse(saved)
		}
	}

	function selectCategory(category) {
		selectedCategory = category
		document.getElementById("categoryTitle").innerText = category + " selected"
	}

	function addReminder() {
		let task = document.getElementById("taskInput").value.trim()
		let dateTime = document.getElementById("dateTimeInput").value
		let priority = document.getElementById("priorityInput").value
		let repeat = document.getElementById("repeatInput").value

		if (selectedCategory === "") {
			alert("Select a category first")
			return
		}

		if (task === "" || dateTime === "") {
			alert("Enter task and date/time")
			return
		}

		let reminder = {
			id: Date.now(),
			category: selectedCategory,
			task: task,
			dateTime: dateTime,
			priority: priority,
			repeat: repeat,
			done: false,
			cancelled: false,
			lastAlertTime: 0,
			createdAt: new Date().toISOString()
		}

		reminders.push(reminder)
		saveReminders()
		renderReminders()

		document.getElementById("taskInput").value = ""
		document.getElementById("dateTimeInput").value = ""
	}

	function formatDateTime(value) {
		let d = new Date(value)
		return d.toLocaleString()
	}

	function getEmoji(category) {
		if (category === "Work") return "💼"
		if (category === "Calls") return "📞"
		if (category === "Medicines") return "💊"
		if (category === "Bills") return "💸"
		if (category === "Personal") return "⭐"
		if (category === "Food") return "🍽️"
		return "🔔"
	}

	function getPriorityColor(priority) {
		if (priority === "Critical") return "#7f1d1d"
		if (priority === "Important") return "#78350f"
		return "#1e293b"
	}

	function renderReminders() {
		let reminderList = document.getElementById("reminderList")
		let search = document.getElementById("searchInput") ? document.getElementById("searchInput").value.toLowerCase() : ""
		let filter = document.getElementById("filterInput") ? document.getElementById("filterInput").value : "All"

		reminderList.innerHTML = ""

		let visible = reminders.filter(function(r) {
			let matchesSearch = r.task.toLowerCase().includes(search) || r.category.toLowerCase().includes(search)

			let matchesFilter =
				filter === "All" ||
				(filter === "Pending" && !r.done && !r.cancelled) ||
				(filter === "Done" && r.done) ||
				(filter === "Cancelled" && r.cancelled)

			return matchesSearch && matchesFilter
		})

		if (visible.length === 0) {
			reminderList.innerHTML = "<div class='empty'>No reminders found.</div>"
		}

		for (let i = 0; i < visible.length; i++) {
			let r = visible[i]

			let card = document.createElement("div")
			card.className = "taskCard"

			let bg = r.done ? "#14532d" : r.cancelled ? "#3f1d1d" : getPriorityColor(r.priority || "Normal")
			let opacity = r.cancelled ? "0.55" : "1"

			card.style.background = bg
			card.style.opacity = opacity

			let statusText = r.done ? "Done" : r.cancelled ? "Cancelled" : "Pending"

			card.innerHTML =
				"<div class='taskTop'>" +
					"<input type='checkbox' " + (r.done ? "checked" : "") + " onchange='toggleDone(" + r.id + ")'>" +
					"<div style='flex:1;'>" +
						"<h3 class='taskTitle'>" + getEmoji(r.category) + " " + r.task + "</h3>" +
						"<p class='meta'><b>Time:</b> " + formatDateTime(r.dateTime) + "</p>" +
						"<span class='badge'>" + r.category + "</span>" +
						"<span class='badge'>" + (r.priority || "Normal") + "</span>" +
						"<span class='badge'>" + (r.repeat || "None") + "</span>" +
						"<span class='badge'>" + statusText + "</span>" +
						"<br>" +
						"<button class='actionBtn' onclick='remindAfter(" + r.id + ", 15)' style='background:#2563eb;'>After 15 min</button>" +
						"<button class='actionBtn' onclick='remindAfter(" + r.id + ", 60)' style='background:#7c3aed;'>After 1 hour</button>" +
						"<button class='actionBtn' onclick='cancelReminder(" + r.id + ")' style='background:#ef4444;'>Cancel</button>" +
						"<button class='actionBtn' onclick='deleteReminder(" + r.id + ")' style='background:#475569;'>Delete</button>" +
					"</div>" +
				"</div>"

			reminderList.appendChild(card)
		}

		updateStats()
	}

	function toggleDone(id) {
		for (let i = 0; i < reminders.length; i++) {
			if (reminders[i].id === id) {
				reminders[i].done = !reminders[i].done
				reminders[i].cancelled = false
			}
		}
		saveReminders()
		renderReminders()
	}

	function cancelReminder(id) {
		for (let i = 0; i < reminders.length; i++) {
			if (reminders[i].id === id) {
				reminders[i].cancelled = true
			}
		}
		saveReminders()
		renderReminders()
	}

	function deleteReminder(id) {
		reminders = reminders.filter(function(r) {
			return r.id !== id
		})
		saveReminders()
		renderReminders()
	}

	function clearCompleted() {
		reminders = reminders.filter(function(r) {
			return !r.done
		})
		saveReminders()
		renderReminders()
	}

	function remindAfter(id, minutes) {
		for (let i = 0; i < reminders.length; i++) {
			if (reminders[i].id === id) {
				let newTime = new Date()
				newTime.setMinutes(newTime.getMinutes() + minutes)
				reminders[i].dateTime = toLocalDateTime(newTime)
				reminders[i].lastAlertTime = 0
				reminders[i].done = false
				reminders[i].cancelled = false
				alert("Okay, moved " + minutes + " minutes later.")
			}
		}
		saveReminders()
		renderReminders()
	}

	function toLocalDateTime(date) {
		return date.getFullYear() + "-" +
			String(date.getMonth() + 1).padStart(2, "0") + "-" +
			String(date.getDate()).padStart(2, "0") + "T" +
			String(date.getHours()).padStart(2, "0") + ":" +
			String(date.getMinutes()).padStart(2, "0")
	}

	function checkReminders() {
		let now = new Date()
		let nowTime = now.getTime()

		for (let i = 0; i < reminders.length; i++) {
			let r = reminders[i]

			if (r.done || r.cancelled) {
				continue
			}

			let reminderTime = new Date(r.dateTime).getTime()

			if (nowTime >= reminderTime) {
				if (r.lastAlertTime === 0 || nowTime - r.lastAlertTime >= 120000) {
					showNotification(r)
					playBeep(r.priority || "Normal")
					r.lastAlertTime = nowTime
					saveReminders()
					renderReminders()
				}
			}
		}
	}

	function showNotification(r) {
		let message = r.category + " - " + r.task + " | " + (r.priority || "Normal")

		if ("Notification" in window && Notification.permission === "granted") {
			window.sendNotification("NUDGE Reminder", message)
		} else {
			alert(message)
		}
	}

	function updateStats() {
		let total = reminders.filter(function(r) {
			return !r.cancelled
		}).length

		let done = reminders.filter(function(r) {
			return r.done && !r.cancelled
		}).length

		let percent = total === 0 ? 0 : Math.round((done / total) * 100)

		document.getElementById("statsText").innerText = done + "/" + total + " tasks done"
		document.getElementById("progressBar").style.width = percent + "%"
	}

	setInterval(checkReminders, 10000)
</script>

<script type="module">
	import { initializeApp } from "https://www.gstatic.com/firebasejs/10.12.2/firebase-app.js";
	import { getMessaging, getToken, onMessage } from "https://www.gstatic.com/firebasejs/10.12.2/firebase-messaging.js";

	const firebaseConfig = {
		apiKey: "AIzaSyCkx84KA8KBzcU3JCrDc-4-LmdyKN6EIL0",
		authDomain: "nudge-5d054.firebaseapp.com",
		projectId: "nudge-5d054",
		storageBucket: "nudge-5d054.firebasestorage.app",
		messagingSenderId: "135553800102",
		appId: "1:135553800102:web:164f023344c8af8da56167",
		measurementId: "G-223C3TPMC0"
	};

	const app = initializeApp(firebaseConfig);
	const messaging = getMessaging(app);

	async function setupFirebasePush() {
		try {
			const permission = await Notification.requestPermission();

			if (permission !== "granted") {
				console.log("Notification permission not granted");
				return;
			}

			const token = await getToken(messaging, {
				vapidKey: "BAkJV6bK48sE7zAqA7LTxBH8GP9QO5RPEZy1YhWw1rIBCx-ShPJD_bNz4ryFfJD547Dxons8nHRc2oIrC6KMjPg"
			});

			console.log("Firebase Push Token:", token);
			localStorage.setItem("firebasePushToken", token);

		} catch (error) {
			console.log("Firebase push error:", error);
		}
	}

	setupFirebasePush();

	onMessage(messaging, function(payload) {
		console.log("Firebase message received:", payload);

		let title = "Nudge Reminder";
		let body = "You have a new reminder";

		if (payload.notification) {
			title = payload.notification.title || title;
			body = payload.notification.body || body;
		}

		window.sendNotification(title, body);
	});
</script>

</body>
</html>
	`)
}

func manifestHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "manifest.json")
}

func serviceWorkerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/javascript")
	http.ServeFile(w, r, "service-worker.js")
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/manifest.json", manifestHandler)
	http.HandleFunc("/service-worker.js", serviceWorkerHandler)

	fmt.Println("Server running on http://localhost:8080")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.ListenAndServe(":"+port, nil)
}

importScripts("https://www.gstatic.com/firebasejs/10.12.2/firebase-app-compat.js");
importScripts("https://www.gstatic.com/firebasejs/10.12.2/firebase-messaging-compat.js");

firebase.initializeApp({
  apiKey: "AIzaSyCkx84KA8KBzcU3JCrDc-4-LmdyKN6EIL0",
  authDomain: "nudge-5d054.firebaseapp.com",
  projectId: "nudge-5d054",
  storageBucket: "nudge-5d054.firebasestorage.app",
  messagingSenderId: "135553800102",
  appId: "1:135553800102:web:164f023344c8af8da56167",
  measurementId: "G-223C3TPMC0"
});

const messaging = firebase.messaging();

messaging.onBackgroundMessage(function(payload) {
  self.registration.showNotification(payload.notification.title, {
    body: payload.notification.body,
    icon: "https://cdn-icons-png.flaticon.com/512/1827/1827392.png"
  });
});
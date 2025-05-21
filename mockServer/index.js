const PORT = 3000;
const express = require("express");
const cors = require("cors");

const app = express();

app.use(cors());

const mockServers = {
  1: {
    serverName: "Для педиков",
    chats: [{ channelTitle: "Гомики пишут", url: 1 }],
    voiceChannels: [
      { channelTitle: "Первая комната", url: "room1" },
      { channelTitle: "Вторая комната", url: "room2" },
    ],
  },
  2: {
    serverName: "Для гомиков",
    chats: [{ channelTitle: "Педики", url: 2 }],
    voiceChannels: [
      { channelTitle: "Большой хуй", url: "room3" },
      { channelTitle: "Вторая хуесссссссссс", url: "room4" },
    ],
  },
};

const mockChats = {
  1: {
    title: "Гомики пишут",
    messages: [
      {
        text: "Я гей",
        timestamp: "2025-05-21T13:21:44.4933535+03:00",
        sender: 1337,
      },
      {
        text: "Я тоже",
        timestamp: "2024-05-21T13:21:44.4933535+03:00",
        sender: 1488,
      },
      {
        text: "Я гей",
        timestamp: "2025-05-21T13:21:44.4933535+03:00",
        sender: 1337,
      },
      {
        text: "Я тоже",
        timestamp: "2024-05-21T13:21:44.4933535+03:00",
        sender: 1488,
      },
      {
        text: "Я гей",
        timestamp: "2025-05-21T13:21:44.4933535+03:00",
        sender: 1337,
      },
      {
        text: "Я тоже",
        timestamp: "2024-05-21T13:21:44.4933535+03:00",
        sender: 1488,
      },
      {
        text: "Я гей",
        timestamp: "2025-05-21T13:21:44.4933535+03:00",
        sender: 1337,
      },
      {
        text: "Я тоже",
        timestamp: "2024-05-21T13:21:44.4933535+03:00",
        sender: 1488,
      },
      {
        text: "Я гей",
        timestamp: "2025-05-21T13:21:44.4933535+03:00",
        sender: 1337,
      },
      {
        text: "Я тоже",
        timestamp: "2024-05-21T13:21:44.4933535+03:00",
        sender: 1488,
      },
      {
        text: "Я гей",
        timestamp: "2025-05-21T13:21:44.4933535+03:00",
        sender: 1337,
      },
      {
        text: "Я тоже",
        timestamp: "2024-02-21T13:21:44.4933535+03:00",
        sender: 1488,
      },
    ],
  },
  2: {
    title: "Педики",
    messages: [
      { text: "Я гей", timestamp: 100, sender: 1337 },
      {
        text: "а я нет пошел нахуй пидор ебаный блядь",
        timestamp: 200,
        sender: 1488,
      },
    ],
  },
};

app.get("/channel-info/:channelId", (req, res) => {
  const channelId = req.params.channelId;

  res.send(mockServers[channelId]);
});

/* app.post('/new-participant', (req, res) => {
  try {
    const body = req.body

    channelParticipants[body.channelId][body.channelViewId][body.userUuid] = body

    res.status(200).end()
  } catch (error) {
    res.status(400).end()
  }
}) */

app.get("/chat-info/:chatId", (req, res) => {
  const chatId = req.params.chatId;

  res.send(mockChats[chatId]);
});

app.listen(PORT, () => {
  console.log(PORT);
});

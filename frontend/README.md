# Voting-Go Frontend

Modern React 19 single-page application for the Voting-Go platform. Built with Vite and TypeScript for high performance and developer productivity.

## Features

- **Dynamic Poll List**: Browse active community polls with real-time metadata.
- **Live Voting View**: Cast votes and see instant feedback via WebSockets.
- **Real-time Results**: Interactive charts showing vote distribution as it happens.
- **Poll Cancellation**: integrated management for poll owners to cancel active polls with confirmation safety.
- **Secure Auth Flow**: Protected routes and session management with JWT.
- **Dark Mode UI**: Professional, accessible dark-themed interface using CSS Modules.

## Tech Stack

- **React 19**: Leveraging the latest features and optimizations.
- **Vite**: Ultra-fast build tool and development server.
- **TypeScript**: Full type safety across the component tree.
- **Lucide React**: Clean, consistent iconography.
- **React Router Dom**: Client-side routing for seamless navigation.

## Installation

1. Navigate to the frontend directory:
   ```bash
   cd frontend
   ```

2. Install dependencies:
   ```bash
   npm install
   ```

3. Configure environment:
   ```bash
   cp .env.example .env
   # Edit VITE_API_URL and VITE_WS_URL to point to your backend
   ```

## Usage

### Development Server
```bash
npm run dev
```

### Production Build
```bash
npm run build
```

### Linting
```bash
npm run lint
```

## Structure

- `/src/components`: UI components organized by feature (PollList, VoteView, ResultsView, etc.).
- `/src/assets`: Static assets and global styles.
- `/public`: Public assets and standard favicon/icons.

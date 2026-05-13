import { useState } from 'react';
import { List, PlusCircle, BarChart2, Vote } from 'lucide-react';
import styles from './App.module.css';
import { PollList } from './components/PollList/PollList';
import { CreatePoll } from './components/CreatePoll/CreatePoll';
import { ResultsView } from './components/ResultsView/ResultsView';
import { VoteView } from './components/VoteView/VoteView';

type Tab = 'polls' | 'create' | 'results' | 'vote';

function App() {
  const [activeTab, setActiveTab] = useState<Tab>('polls');
  const [selectedPollId, setSelectedPollId] = useState<string | null>(null);

  const handleSelectPoll = (id: string, view: 'vote' | 'results') => {
    setSelectedPollId(id);
    setActiveTab(view);
  };

  const renderContent = () => {
    switch (activeTab) {
      case 'polls':
        return <PollList onSelectPoll={handleSelectPoll} />;
      case 'create':
        return <CreatePoll onPollCreated={() => setActiveTab('polls')} />;
      case 'vote':
        return selectedPollId ? (
          <VoteView 
            pollId={selectedPollId} 
            onVoteSuccess={() => setActiveTab('results')}
            onBack={() => setActiveTab('polls')}
          />
        ) : <PollList onSelectPoll={handleSelectPoll} />;
      case 'results':
        return selectedPollId ? (
          <ResultsView pollId={selectedPollId} />
        ) : <PollList onSelectPoll={handleSelectPoll} />;
      default:
        return <PollList onSelectPoll={handleSelectPoll} />;
    }
  };

  return (
    <div className={styles.appContainer}>
      <header className={styles.header}>
        <h1>Voting Go</h1>
        <p>Real-time polling made simple</p>
      </header>

      <nav className={styles.tabNav}>
        <button 
          className={`${styles.tabButton} ${activeTab === 'polls' ? styles.active : ''}`}
          onClick={() => setActiveTab('polls')}
        >
          <List size={20} />
          Polls
        </button>
        <button 
          className={`${styles.tabButton} ${activeTab === 'create' ? styles.active : ''}`}
          onClick={() => setActiveTab('create')}
        >
          <PlusCircle size={20} />
          Create Poll
        </button>
        {(activeTab === 'vote' || activeTab === 'results') && selectedPollId && (
          <button className={`${styles.tabButton} ${styles.active}`}>
            {activeTab === 'vote' ? <Vote size={20} /> : <BarChart2 size={20} />}
            {activeTab === 'vote' ? 'Voting' : 'Live Results'}
          </button>
        )}
      </nav>

      <main className="content">
        {renderContent()}
      </main>
    </div>
  );
}

export default App;

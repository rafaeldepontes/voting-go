import { useEffect, useState } from 'react';
import { Vote, BarChart3, RefreshCw } from 'lucide-react';
import styles from './PollList.module.css';

interface PollDto {
  id: string;
  text: string;
}

interface PollListProps {
  onSelectPoll: (id: string, view: 'vote' | 'results') => void;
}

const API_URL = import.meta.env.VITE_API_URL;

export const PollList = ({ onSelectPoll }: PollListProps) => {
  const [polls, setPolls] = useState<PollDto[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchPolls = async () => {
    setLoading(true);
    try {
      const response = await fetch(`${API_URL}/polls`);
      if (!response.ok) throw new Error('Failed to fetch polls');
      const data = await response.json();
      setPolls(data || []);
      setError(null);
    } catch (err) {
      setError('Could not load polls. Is the backend running?');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchPolls();
  }, []);

  return (
    <div className={styles.pollList}>
      <div className={styles.listHeader}>
        <h2>Available Polls</h2>
        <button onClick={fetchPolls} className={styles.iconButton} title="Refresh">
          <RefreshCw size={20} className={loading ? styles.spinning : ''} />
        </button>
      </div>

      {loading && <p>Loading polls...</p>}
      {error && <p className={styles.error}>{error}</p>}
      {!loading && !error && polls.length === 0 && (
        <p>No polls created yet. Be the first!</p>
      )}

      <div className={styles.pollGrid}>
        {polls.map((poll) => (
          <div key={poll.id} className={styles.pollCard}>
            <h3>{poll.text}</h3>
            <div className={styles.cardActions}>
              <button 
                onClick={() => onSelectPoll(poll.id, 'vote')}
                className={`${styles.actionButton} ${styles.vote}`}
              >
                <Vote size={18} />
                Vote
              </button>
              <button 
                onClick={() => onSelectPoll(poll.id, 'results')}
                className={`${styles.actionButton} ${styles.results}`}
              >
                <BarChart3 size={18} />
                Results
              </button>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

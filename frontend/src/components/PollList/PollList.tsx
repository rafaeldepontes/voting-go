import { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { Vote, BarChart2, PlusCircle, Clock, AlertCircle, Loader2 } from 'lucide-react';
import styles from './PollList.module.css';

interface Poll {
  id: string;
  text: string;
  duration: number;
}

interface PollListProps {
  token: string | null;
  onAuthError: () => void;
}

const API_URL = import.meta.env.VITE_API_URL;

export const PollList = ({ token, onAuthError }: PollListProps) => {
  const [polls, setPolls] = useState<Poll[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetchPolls();
  }, [token]);

  const fetchPolls = async () => {
    setLoading(true);
    setError(null);
    try {
      const headers: Record<string, string> = {};
      if (token) {
        headers['Authorization'] = `Bearer ${token}`;
      }

      const response = await fetch(`${API_URL}/polls`, { headers });

      if (response.status === 401 || response.status === 403) {
        if (token) onAuthError();
      }

      if (!response.ok) {
        throw new Error('Failed to fetch polls');
      }

      const data = await response.json();
      setPolls(data || []);
    } catch (err: any) {
      setError(err.message || 'Could not load polls');
    } finally {
      setLoading(false);
    }
  };

  const formatDuration = (nanoseconds: number) => {
    if (nanoseconds <= 0) return 'Permanent';
    const minutes = Math.floor(nanoseconds / 60000000000);
    const hours = Math.floor(minutes / 60);
    if (hours > 0) return `${hours}h ${minutes % 60}m`;
    return `${minutes}m`;
  };

  if (loading) {
    return (
      <div className={styles.loadingState}>
        <Loader2 size={40} className={styles.spinner} />
        <p>Loading active polls...</p>
      </div>
    );
  }

  return (
    <div className={styles.pollListPage}>
      <header className={styles.pageHeader}>
        <div className={styles.headerTitle}>
          <h1>Browse Polls</h1>
          <p>Participate in real-time community decisions</p>
        </div>
        <Link to="/create" className={styles.createButton}>
          <PlusCircle size={20} />
          <span>Create New Poll</span>
        </Link>
      </header>

      {error && (
        <div className={styles.errorState}>
          <AlertCircle size={24} />
          <p>{error}</p>
          <button onClick={fetchPolls} className={styles.retryButton}>Try Again</button>
        </div>
      )}

      {!error && polls.length === 0 ? (
        <div className={styles.emptyState}>
          <Vote size={48} className={styles.emptyIcon} />
          <h2>No polls found</h2>
          <p>Be the first one to create a poll!</p>
          <Link to="/create" className={styles.emptyLink}>Create Poll</Link>
        </div>
      ) : (
        <div className={styles.pollGrid}>
          {polls.map((poll) => (
            <div key={poll.id} className={styles.pollCard}>
              <div className={styles.pollInfo}>
                <h3 className={styles.pollText}>{poll.text}</h3>
                <div className={styles.pollMeta}>
                  <div className={styles.metaItem}>
                    <Clock size={14} />
                    <span>{formatDuration(poll.duration)}</span>
                  </div>
                </div>
              </div>
              <div className={styles.pollActions}>
                <Link to={`/poll/${poll.id}`} className={styles.voteLink}>
                  <Vote size={18} />
                  <span>Vote</span>
                </Link>
                <Link to={`/results/${poll.id}`} className={styles.resultsLink}>
                  <BarChart2 size={18} />
                  <span>Results</span>
                </Link>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

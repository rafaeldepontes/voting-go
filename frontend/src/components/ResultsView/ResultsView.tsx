import { useState, useEffect } from 'react';
import { useParams, useNavigate, Link } from 'react-router-dom';
import { BarChart2, ArrowLeft, Loader2, Users, Clock, Vote, Trash2 } from 'lucide-react';
import styles from './ResultsView.module.css';

interface Option {
  id: number;
  text: string;
  votes: number;
}

interface Poll {
  id: string;
  text: string;
  options: Option[];
  duration: number;
  createdAt: string;
}

interface ResultsViewProps {
  token: string | null;
  onAuthError?: () => void;
}

const API_URL = import.meta.env.VITE_API_URL;
const WS_URL = import.meta.env.VITE_WS_URL || API_URL?.replace('http', 'ws');

export const ResultsView = ({ token, onAuthError }: ResultsViewProps) => {
  const { id: pollId } = useParams<{ id: string }>();
  const [poll, setPoll] = useState<Poll | null>(null);
  const [loading, setLoading] = useState(true);
  const [cancelling, setCancelling] = useState(false);
  const [showConfirm, setShowConfirm] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const navigate = useNavigate();

  useEffect(() => {
    if (!pollId) return;

    let socket: WebSocket | null = null;
    const connectWS = () => {
      const tokenQuery = token ? `?token=${token}` : '';
      socket = new WebSocket(`${WS_URL}/ws/polls/${pollId}${tokenQuery}`);

      socket.onmessage = (event) => {
        const updatedPoll = JSON.parse(event.data);
        setPoll(updatedPoll);
        setLoading(false);
      };

      socket.onerror = () => {
        console.error('Real-time connection failed.');
        setLoading(false);
      };
    };

    connectWS();

    return () => {
      if (socket) socket.close();
    };
  }, [pollId, token]);

  const handleCancelPoll = async () => {
    if (!pollId) return;

    setCancelling(true);
    setError(null);

    try {
      const headers: Record<string, string> = {};
      if (token) {
        headers['Authorization'] = `Bearer ${token}`;
      }

      const response = await fetch(`${API_URL}/polls/${pollId}`, {
        method: 'DELETE',
        headers,
      });

      if (response.status === 401 || response.status === 403) {
        if (token && onAuthError) {
          onAuthError();
          return;
        }
        throw new Error('You do not have permission to cancel this poll');
      }

      if (!response.ok) {
        const errorText = await response.text();
        throw new Error(errorText || 'Failed to cancel poll');
      }

      navigate('/');
    } catch (err: any) {
      setError(err.message || 'Could not cancel poll.');
      setShowConfirm(false);
    } finally {
      setCancelling(false);
    }
  };

  if (loading) {
    return (
      <div className={styles.loadingState}>
        <Loader2 size={40} className={styles.spinner} />
        <p>Connecting to live results...</p>
      </div>
    );
  }

  if (!poll) {
    return (
      <div className={styles.errorState}>
        <BarChart2 size={40} />
        <h2>Poll results unavailable</h2>
        <p>We couldn't find the results for this poll.</p>
        <Link to="/" className={styles.backLink}>Back to Polls</Link>
      </div>
    );
  }

  const totalVotes = poll.options.reduce((sum, opt) => sum + opt.votes, 0);

  const getTimeRemaining = () => {
    if (poll.duration <= 0) return 'Permanent';
    const createdAt = new Date(poll.createdAt).getTime();
    const durationMs = poll.duration / 1000000;
    const expiresAt = createdAt + durationMs;
    const remaining = expiresAt - Date.now();

    if (remaining <= 0) return 'Concluded';

    const minutes = Math.floor(remaining / 60000);
    const hours = Math.floor(minutes / 60);
    if (hours > 0) return `${hours}h ${minutes % 60}m remaining`;
    return `${minutes}m remaining`;
  };

  return (
    <div className={styles.resultsPage}>
      <div className={styles.topActions}>
        <button onClick={() => navigate(-1)} className={styles.backButton}>
          <ArrowLeft size={18} />
          <span>Back</span>
        </button>
      </div>

      <div className={styles.resultsCard}>
        <header className={styles.resultsHeader}>
          <div className={styles.headerTop}>
            <div className={styles.pollMeta}>
              <div className={styles.metaItem}>
                <Clock size={14} />
                <span>{getTimeRemaining()}</span>
              </div>
              <div className={styles.metaItem}>
                <Users size={14} />
                <span>{totalVotes} total votes</span>
              </div>
            </div>
            {token && (
              <button
                className={styles.deleteButtonInside}
                onClick={() => setShowConfirm(true)}
                title="Cancel poll"
              >
                <Trash2 size={18} />
              </button>
            )}
          </div>
          <h1>{poll.text}</h1>
          <p>Live results updated in real-time as votes are cast.</p>
        </header>

        {error && (
          <div className={styles.errorBanner}>
            <Trash2 size={18} />
            <span>{error}</span>
          </div>
        )}

        <div className={styles.optionsGrid}>
          {poll.options.sort((a, b) => b.votes - a.votes).map((option) => {
            const percentage = totalVotes > 0 ? (option.votes / totalVotes) * 100 : 0;
            return (
              <div key={option.id} className={styles.optionResult}>
                <div className={styles.optionInfo}>
                  <span className={styles.optionText}>{option.text}</span>
                  <span className={styles.voteCount}>
                    {option.votes} {option.votes === 1 ? 'vote' : 'votes'} ({percentage.toFixed(1)}%)
                  </span>
                </div>
                <div className={styles.progressContainer}>
                  <div
                    className={styles.progressBar}
                    style={{ width: `${percentage}%` }}
                  ></div>
                </div>
              </div>
            );
          })}
        </div>

        <div className={styles.footerActions}>
          <Link to={`/poll/${pollId}`} className={styles.voteLink}>
            <Vote size={18} />
            <span>Cast Vote</span>
          </Link>
        </div>
      </div>

      {showConfirm && (
        <div className={styles.modalOverlay}>
          <div className={styles.modal}>
            <h2>Cancel Poll</h2>
            <p>Are you sure you want to cancel this poll? This action cannot be undone and all results will be lost.</p>
            <div className={styles.modalActions}>
              <button
                className={styles.confirmButton}
                onClick={handleCancelPoll}
                disabled={cancelling}
              >
                {cancelling ? <Loader2 size={18} className={styles.spinner} /> : 'Yes, Cancel Poll'}
              </button>
              <button
                className={styles.cancelButton}
                onClick={() => setShowConfirm(false)}
                disabled={cancelling}
              >
                No, Keep It
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};


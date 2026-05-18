import { useState, useEffect } from 'react';
import { useParams, useNavigate, Link } from 'react-router-dom';
import { Vote, ArrowLeft, AlertCircle, Loader2, CheckCircle, Clock, Users, Trash2 } from 'lucide-react';
import styles from './VoteView.module.css';

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

interface VoteViewProps {
  token: string | null;
  onAuthError: () => void;
  onVoteSuccess: () => void;
}

const API_URL = import.meta.env.VITE_API_URL;
const WS_URL = import.meta.env.VITE_WS_URL || API_URL?.replace('http', 'ws');

export const VoteView = ({ token, onAuthError, onVoteSuccess }: VoteViewProps) => {
  const { id: pollId } = useParams<{ id: string }>();
  const [poll, setPoll] = useState<Poll | null>(null);
  const [selectedOption, setSelectedOption] = useState<number | null>(null);
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [cancelling, setCancelling] = useState(false);
  const [showConfirm, setShowConfirm] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [voted, setVoted] = useState(false);
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
        setError('Real-time connection failed. Falling back to list page.');
        setLoading(false);
      };

      socket.onclose = () => {
        console.log('WS connection closed');
      };
    };

    connectWS();

    return () => {
      if (socket) socket.close();
    };
  }, [pollId, token]);

  const handleVote = async () => {
    if (selectedOption === null || !pollId) return;

    setSubmitting(true);
    setError(null);

    try {
      const headers: Record<string, string> = { 'Content-Type': 'application/json' };
      if (token) {
        headers['Authorization'] = `Bearer ${token}`;
      }

      const response = await fetch(`${API_URL}/polls/${pollId}/vote`, {
        method: 'POST',
        headers,
        body: JSON.stringify({ optionId: selectedOption }),
      });

      if (response.status === 401 || response.status === 403) {
        if (token) onAuthError();
      }

      if (!response.ok) {
        const errorText = await response.text();
        throw new Error(errorText || 'Failed to register vote');
      }

      setVoted(true);
      onVoteSuccess();
      setTimeout(() => {
        navigate(`/results/${pollId}`);
      }, 1500);
    } catch (err: any) {
      setError(err.message || 'Could not register your vote.');
    } finally {
      setSubmitting(false);
    }
  };

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
        if (token) {
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

  const getTimeRemaining = () => {
    if (!poll || poll.duration <= 0) return 'Permanent';
    const createdAt = new Date(poll.createdAt).getTime();
    const durationMs = poll.duration / 1000000;
    const expiresAt = createdAt + durationMs;
    const remaining = expiresAt - Date.now();
    
    if (remaining <= 0) return 'Expired';
    
    const minutes = Math.floor(remaining / 60000);
    const hours = Math.floor(minutes / 60);
    if (hours > 0) return `${hours}h ${minutes % 60}m remaining`;
    return `${minutes}m remaining`;
  };

  if (loading) {
    return (
      <div className={styles.loadingState}>
        <Loader2 size={40} className={styles.spinner} />
        <p>Connecting to live poll...</p>
      </div>
    );
  }

  if (!poll) {
    return (
      <div className={styles.errorState}>
        <AlertCircle size={40} />
        <h2>Poll not found</h2>
        <p>The poll you are looking for doesn't exist or has been removed.</p>
        <Link to="/" className={styles.backLink}>Back to Polls</Link>
      </div>
    );
  }

  const isExpired = getTimeRemaining() === 'Expired';
  const totalVotes = poll.options.reduce((sum, opt) => sum + opt.votes, 0);

  return (
    <div className={styles.votePage}>
      <div className={styles.topActions}>
        <button onClick={() => navigate('/')} className={styles.backButton}>
          <ArrowLeft size={18} />
          <span>Back to Polls</span>
        </button>
      </div>

      <div className={styles.voteCard}>
        <header className={styles.pollHeader}>
          <div className={styles.headerTop}>
            <div className={styles.pollMeta}>
              <div className={styles.metaItem}>
                <Clock size={14} />
                <span>{getTimeRemaining()}</span>
              </div>
              <div className={styles.metaItem}>
                <Users size={14} />
                <span>{totalVotes} votes cast</span>
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
          <p>Select one of the options below to cast your vote.</p>
        </header>

        {voted ? (
          <div className={styles.successState}>
            <CheckCircle size={48} className={styles.successIcon} />
            <h2>Vote Registered!</h2>
            <p>Thank you for participating. Redirecting to results...</p>
          </div>
        ) : (
          <div className={styles.optionsList}>
            {poll.options.map((option) => (
              <button
                key={option.id}
                className={`${styles.optionButton} ${selectedOption === option.id ? styles.selected : ''}`}
                onClick={() => !isExpired && setSelectedOption(option.id)}
                disabled={submitting || isExpired}
              >
                <span className={styles.optionText}>{option.text}</span>
                {selectedOption === option.id && <CheckCircle size={20} className={styles.checkIcon} />}
              </button>
            ))}
          </div>
        )}

        {error && (
          <div className={styles.errorBanner}>
            <AlertCircle size={18} />
            <span>{error}</span>
          </div>
        )}

        {!voted && (
          <div className={styles.actions}>
            <button
              className={styles.submitButton}
              onClick={handleVote}
              disabled={selectedOption === null || submitting || isExpired}
            >
              {submitting ? <Loader2 size={18} className={styles.spinner} /> : (
                <>
                  <Vote size={18} />
                  <span>{isExpired ? 'Poll Expired' : 'Cast Your Vote'}</span>
                </>
              )}
            </button>
            <Link to={`/results/${pollId}`} className={styles.resultsLink}>
              View Live Results
            </Link>
          </div>
        )}
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


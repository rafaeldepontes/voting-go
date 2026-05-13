import { useEffect, useState } from 'react';
import { CheckCircle2, ArrowLeft } from 'lucide-react';
import styles from './VoteView.module.css';

interface Option {
  id: number;
  text: string;
}

interface Poll {
  id: string;
  text: string;
  options: Option[];
}

interface VoteViewProps {
  pollId: string;
  onVoteSuccess: () => void;
  onBack: () => void;
}

const API_URL = import.meta.env.VITE_API_URL;
const WS_URL = import.meta.env.VITE_WS_URL;

export const VoteView = ({ pollId, onVoteSuccess, onBack }: VoteViewProps) => {
  const [poll, setPoll] = useState<Poll | null>(null);
  const [selectedOption, setSelectedOption] = useState<number | null>(null);
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const ws = new WebSocket(`${WS_URL}/ws/polls/${pollId}`);
    ws.onmessage = (event) => {
      const data = JSON.parse(event.data);
      setPoll(data);
      setLoading(false);
      ws.close();
    };
    ws.onerror = () => {
      setError('Failed to load poll options.');
      setLoading(false);
    };
  }, [pollId]);

  const handleVote = async () => {
    if (selectedOption === null) return;

    setSubmitting(true);
    setError(null);

    try {
      const response = await fetch(`${API_URL}/polls/${pollId}/vote`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ optionId: selectedOption }),
      });

      if (!response.ok) throw new Error('Failed to register vote');
      
      onVoteSuccess();
    } catch (err) {
      setError('Could not submit vote. Is the backend running?');
    } finally {
      setSubmitting(false);
    }
  };

  if (loading) return <p>Loading options...</p>;
  if (error) return <p className={styles.error}>{error}</p>;
  if (!poll) return <p>Poll not found.</p>;

  return (
    <div className={styles.voteView}>
      <button onClick={onBack} className={styles.backLink}>
        <ArrowLeft size={16} /> Back to polls
      </button>
      
      <h2>{poll.text}</h2>
      <p className="instruction">Select an option below:</p>

      <div className={styles.optionsList}>
        {poll.options.map((option) => (
          <label 
            key={option.id} 
            className={`${styles.optionItem} ${selectedOption === option.id ? styles.selected : ''}`}
          >
            <input
              type="radio"
              name="poll-option"
              value={option.id}
              checked={selectedOption === option.id}
              onChange={() => setSelectedOption(option.id)}
              disabled={submitting}
            />
            <span className={styles.optionLabelText}>{option.text}</span>
            {selectedOption === option.id && <CheckCircle2 size={20} className={styles.checkIcon} />}
          </label>
        ))}
      </div>

      <button 
        onClick={handleVote} 
        className={styles.submitVoteButton} 
        disabled={selectedOption === null || submitting}
      >
        {submitting ? 'Voting...' : 'Submit Vote'}
      </button>
    </div>
  );
};

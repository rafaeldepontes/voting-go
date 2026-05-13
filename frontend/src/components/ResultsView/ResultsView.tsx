import { useEffect, useState } from 'react';
import { Users, Info } from 'lucide-react';
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
}

interface ResultsViewProps {
  pollId: string;
}

const WS_URL = import.meta.env.VITE_WS_URL;

export const ResultsView = ({ pollId }: ResultsViewProps) => {
  const [poll, setPoll] = useState<Poll | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [connected, setConnected] = useState(false);

  useEffect(() => {
    const ws = new WebSocket(`${WS_URL}/ws/polls/${pollId}`);

    ws.onopen = () => {
      setConnected(true);
      setError(null);
    };

    ws.onmessage = (event) => {
      const data = JSON.parse(event.data);
      setPoll(data);
    };

    ws.onerror = () => {
      setError('WebSocket connection error. Is the backend running?');
      setConnected(false);
    };

    ws.onclose = () => {
      setConnected(false);
    };

    return () => {
      ws.close();
    };
  }, [pollId]);

  if (error) return <p className={styles.error}>{error}</p>;
  if (!poll) return <p>Connecting to live results...</p>;

  const totalVotes = poll.options.reduce((sum, opt) => sum + opt.votes, 0);

  return (
    <div className={styles.resultsView}>
      <div className={styles.resultsHeader}>
        <h2>{poll.text}</h2>
        <div className={`${styles.statusBadge} ${connected ? styles.live : styles.offline}`}>
          {connected ? '● LIVE' : '○ OFFLINE'}
        </div>
      </div>

      <div className={styles.totalVotes}>
        <Users size={20} />
        <span>{totalVotes} total votes</span>
      </div>

      <div className={styles.optionsResults}>
        {poll.options.map((option) => {
          const percentage = totalVotes === 0 ? 0 : Math.round((option.votes / totalVotes) * 100);
          return (
            <div key={option.id} className={styles.resultRow}>
              <div className={styles.resultInfo}>
                <span className={styles.optionText}>{option.text}</span>
                <span className={styles.voteCount}>{option.votes} votes ({percentage}%)</span>
              </div>
              <div className={styles.progressBarBg}>
                <div 
                  className={styles.progressBarFill} 
                  style={{ width: `${percentage}%` }}
                ></div>
              </div>
            </div>
          );
        })}
      </div>

      <div className={styles.liveHint}>
        <Info size={16} />
        <span>Results update automatically as votes come in.</span>
      </div>
    </div>
  );
};

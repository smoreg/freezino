import { useState, useEffect, useRef } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { useTranslation } from 'react-i18next';
import axios from 'axios';
import { WinConfetti } from '../animations';
import { useSound } from '../../hooks/useSound';

// Slot symbols
const SYMBOLS = ['🍒', '🍋', '🍊', '🍇', '💎', '⭐', '7️⃣'];

// Payout table
const PAYOUTS = {
  '7️⃣': { 5: 500, 4: 100, 3: 20 },
  '⭐': { 5: 200, 4: 50, 3: 10 },
  '💎': { 5: 150, 4: 40, 3: 8 },
  '🍇': { 5: 100, 4: 25, 3: 5 },
  '🍊': { 5: 80, 4: 20, 3: 4 },
  '🍋': { 5: 60, 4: 15, 3: 3 },
  '🍒': { 5: 40, 4: 10, 3: 2 },
};

// Bet options
const BET_OPTIONS = [10, 25, 50, 100, 250, 500];

interface WinningLine {
  line_number: number;
  symbol: string;
  count: number;
  multiplier: number;
  win: number;
}

interface SlotResult {
  reels: string[][];
  winning_line: WinningLine[];
  total_win: number;
  multiplier: number;
}

interface SlotsProps {
  userBalance: number;
  userId: number;
  onBalanceChange?: (newBalance: number) => void;
}

const Slots = ({ userBalance, userId, onBalanceChange }: SlotsProps) => {
  const { t } = useTranslation();
  const { playSound } = useSound();
  const [reels, setReels] = useState<string[][]>([
    ['🍒', '🍋', '🍊'],
    ['🍇', '💎', '⭐'],
    ['7️⃣', '🍒', '🍋'],
    ['🍊', '🍇', '💎'],
    ['⭐', '7️⃣', '🍒'],
  ]);
  const [spinning, setSpinning] = useState(false);
  const [selectedBet, setSelectedBet] = useState(BET_OPTIONS[0]);
  const [winAmount, setWinAmount] = useState<number>(0);
  const [winningLines, setWinningLines] = useState<WinningLine[]>([]);
  const [showPaytable, setShowPaytable] = useState(false);
  const [balance, setBalance] = useState(userBalance);
  const [message, setMessage] = useState('');
  const [showConfetti, setShowConfetti] = useState(false);
  const spinTimeouts = useRef<ReturnType<typeof setTimeout>[]>([]);

  useEffect(() => {
    setBalance(userBalance);
  }, [userBalance]);

  // Generate random symbols for animation
  const generateRandomSymbols = () => {
    return Array(3)
      .fill(0)
      .map(() => SYMBOLS[Math.floor(Math.random() * SYMBOLS.length)]);
  };

  // Spin animation
  const animateReels = (finalReels: string[][]) => {
    setSpinning(true);
    setWinAmount(0);
    setWinningLines([]);
    setMessage('');

    // Play spin sound
    playSound('slot-spin', 0.5);

    // Clear any existing timeouts
    spinTimeouts.current.forEach(timeout => clearTimeout(timeout));
    spinTimeouts.current = [];

    // Animate each reel
    reels.forEach((_, index) => {
      let spinCount = 0;
      const maxSpins = 20 + index * 5; // Each reel spins longer

      const spinInterval = setInterval(() => {
        setReels(prevReels => {
          const newReels = [...prevReels];
          newReels[index] = generateRandomSymbols();
          return newReels;
        });

        spinCount++;

        if (spinCount >= maxSpins) {
          clearInterval(spinInterval);
          setReels(prevReels => {
            const newReels = [...prevReels];
            newReels[index] = finalReels[index];
            return newReels;
          });

          // Play stop sound for each reel
          playSound('slot-stop', 0.4);

          // If this is the last reel, stop spinning
          if (index === reels.length - 1) {
            const timeout = setTimeout(() => {
              setSpinning(false);
            }, 300);
            spinTimeouts.current.push(timeout);
          }
        }
      }, 50);
    });
  };

  // Handle spin
  const handleSpin = async () => {
    if (spinning) return;
    if (balance < selectedBet) {
      setMessage(t('slots.insufficientBalance') || 'Insufficient balance!');
      return;
    }

    try {
      const response = await axios.post(
        `/api/games/slots/spin?user_id=${userId}`,
        { bet: selectedBet }
      );

      if (response.data.success) {
        const result: SlotResult = response.data.data.result;
        const newBalance = response.data.data.new_balance;

        // Animate reels with final result
        animateReels(result.reels);

        // Update balance and results after animation
        setTimeout(() => {
          setBalance(newBalance);
          setWinAmount(result.total_win);
          setWinningLines(result.winning_line || []);

          if (onBalanceChange) {
            onBalanceChange(newBalance);
          }

          // Play win or lose sound
          if (result.total_win > 0) {
            playSound('win', 0.6);
            setMessage(
              t('slots.youWon', { amount: result.total_win.toFixed(2) }) ||
                `You won $${result.total_win.toFixed(2)}!`
            );
            setShowConfetti(true);
          } else {
            playSound('lose', 0.4);
            setMessage(t('slots.tryAgain') || 'Try again!');
          }
        }, (reels.length + 1) * 300);
      }
    } catch (error: unknown) {
      console.error('Spin error:', error);
      setMessage(
        (error as { response?: { data?: { message?: string } } }).response?.data?.message ||
          t('slots.error') ||
          'An error occurred'
      );
      setSpinning(false);
    }
  };

  // Cleanup timeouts on unmount
  useEffect(() => {
    return () => {
      spinTimeouts.current.forEach(timeout => clearTimeout(timeout));
    };
  }, []);

  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-900 via-purple-900 to-gray-900 p-8">
      <div className="max-w-6xl mx-auto">
        {/* Header */}
        <motion.div
          initial={{ opacity: 0, y: -20 }}
          animate={{ opacity: 1, y: 0 }}
          className="text-center mb-8"
        >
          <h1 className="text-5xl font-bold text-transparent bg-clip-text bg-gradient-to-r from-yellow-400 via-red-500 to-pink-500 mb-4">
            🎰 {t('slots.title') || 'Slot Machine'}
          </h1>
          <div className="flex justify-center items-center gap-8 text-white">
            <div className="text-xl">
              <span className="text-gray-400">{t('slots.balance') || 'Balance'}:</span>{' '}
              <span className="font-bold text-green-400">${balance.toFixed(2)}</span>
            </div>
            <button
              onClick={() => setShowPaytable(!showPaytable)}
              className="px-4 py-2 bg-purple-600 hover:bg-purple-700 rounded-lg transition-colors"
            >
              {showPaytable
                ? t('slots.hidePaytable') || 'Hide Paytable'
                : t('slots.showPaytable') || 'Show Paytable'}
            </button>
          </div>
        </motion.div>

        {/* Paytable */}
        <AnimatePresence>
          {showPaytable && (
            <motion.div
              initial={{ opacity: 0, height: 0 }}
              animate={{ opacity: 1, height: 'auto' }}
              exit={{ opacity: 0, height: 0 }}
              className="mb-8 bg-gray-800/50 backdrop-blur-sm rounded-xl p-6 border border-purple-500/30"
            >
              <h2 className="text-2xl font-bold text-yellow-400 mb-4 text-center">
                {t('slots.paytable') || 'Paytable'}
              </h2>
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                {Object.entries(PAYOUTS).map(([symbol, payouts]) => (
                  <div
                    key={symbol}
                    className="bg-gray-700/50 rounded-lg p-4 border border-gray-600"
                  >
                    <div className="text-4xl text-center mb-2">{symbol}</div>
                    <div className="text-white text-sm space-y-1">
                      <div className="flex justify-between">
                        <span>5x:</span>
                        <span className="text-yellow-400 font-bold">{payouts[5]}x</span>
                      </div>
                      <div className="flex justify-between">
                        <span>4x:</span>
                        <span className="text-yellow-400 font-bold">{payouts[4]}x</span>
                      </div>
                      <div className="flex justify-between">
                        <span>3x:</span>
                        <span className="text-yellow-400 font-bold">{payouts[3]}x</span>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
              <div className="mt-4 text-center text-gray-400 text-sm">
                {t('slots.paytableInfo') ||
                  'Win by matching 3+ symbols on any of the 10 paylines'}
              </div>
            </motion.div>
          )}
        </AnimatePresence>

        {/* Slot Machine */}
        <motion.div
          initial={{ opacity: 0, scale: 0.9 }}
          animate={{ opacity: 1, scale: 1 }}
          className="bg-gradient-to-b from-red-900 to-red-700 rounded-3xl p-8 shadow-2xl border-8 border-yellow-600"
        >
          {/* Reels */}
          <div className="bg-black/50 rounded-2xl p-6 mb-6">
            <div className="flex justify-center gap-4">
              {reels.map((reel, reelIndex) => (
                <div
                  key={reelIndex}
                  className="bg-white/10 rounded-xl p-2 border-4 border-yellow-500/50"
                >
                  <div className="flex flex-col gap-2">
                    {reel.map((symbol, symbolIndex) => (
                      <motion.div
                        key={`${reelIndex}-${symbolIndex}`}
                        animate={spinning ? { y: [0, -10, 0] } : { y: 0 }}
                        transition={{
                          duration: 0.1,
                          repeat: spinning ? Infinity : 0,
                        }}
                        className={`text-6xl w-20 h-20 flex items-center justify-center ${
                          symbolIndex === 1 ? 'bg-yellow-500/20 rounded-lg' : ''
                        }`}
                      >
                        {symbol}
                      </motion.div>
                    ))}
                  </div>
                </div>
              ))}
            </div>
          </div>

          {/* Win Display */}
          <AnimatePresence>
            {winAmount > 0 && (
              <motion.div
                initial={{ opacity: 0, scale: 0.5 }}
                animate={{ opacity: 1, scale: 1 }}
                exit={{ opacity: 0, scale: 0.5 }}
                className="text-center mb-4"
              >
                <div className="text-4xl font-bold text-yellow-400 mb-2">
                  🎉 {t('slots.win') || 'WIN'}: ${winAmount.toFixed(2)} 🎉
                </div>
                {winningLines.map((line, index) => (
                  <div key={index} className="text-white text-sm">
                    {t('slots.line') || 'Line'} {line.line_number}: {line.symbol} x{line.count} ={' '}
                    {line.multiplier}x (${line.win.toFixed(2)})
                  </div>
                ))}
              </motion.div>
            )}
          </AnimatePresence>

          {/* Message */}
          {message && (
            <div
              className={`text-center mb-4 text-lg font-semibold ${
                winAmount > 0 ? 'text-green-400' : 'text-red-400'
              }`}
            >
              {message}
            </div>
          )}

          {/* Bet Selection */}
          <div className="mb-6">
            <div className="text-white text-center mb-3 text-lg font-semibold">
              {t('slots.selectBet') || 'Select Bet'}
            </div>
            <div className="flex justify-center gap-2 flex-wrap">
              {BET_OPTIONS.map(bet => (
                <button
                  key={bet}
                  onClick={() => setSelectedBet(bet)}
                  disabled={spinning}
                  className={`px-6 py-3 rounded-lg font-bold transition-all ${
                    selectedBet === bet
                      ? 'bg-yellow-500 text-black scale-110'
                      : 'bg-gray-700 text-white hover:bg-gray-600'
                  } ${spinning ? 'opacity-50 cursor-not-allowed' : ''}`}
                >
                  ${bet}
                </button>
              ))}
            </div>
          </div>

          {/* Spin Button */}
          <div className="text-center">
            <motion.button
              whileHover={!spinning ? { scale: 1.05 } : {}}
              whileTap={!spinning ? { scale: 0.95 } : {}}
              onClick={handleSpin}
              disabled={spinning || balance < selectedBet}
              className={`px-16 py-6 rounded-2xl text-3xl font-bold transition-all shadow-lg ${
                spinning || balance < selectedBet
                  ? 'bg-gray-600 text-gray-400 cursor-not-allowed'
                  : 'bg-gradient-to-r from-green-500 to-emerald-600 text-white hover:shadow-green-500/50'
              }`}
            >
              {spinning ? '🎰 ' + (t('slots.spinning') || 'SPINNING...') : '🎰 ' + (t('slots.spin') || 'SPIN')}
            </motion.button>
          </div>
        </motion.div>

        {/* Info */}
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.3 }}
          className="mt-8 text-center text-gray-400 text-sm"
        >
          <p>{t('slots.info') || '10 paylines • Match 3+ symbols to win • Higher bets = higher wins'}</p>
        </motion.div>
      </div>

      {/* Win Confetti */}
      <WinConfetti active={showConfetti} onComplete={() => setShowConfetti(false)} />
    </div>
  );
};

export default Slots;

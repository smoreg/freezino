import { useState, type FormEvent } from 'react';
import { motion } from 'framer-motion';
import axios from 'axios';

interface ContactFormData {
  name: string;
  email: string;
  message: string;
}

interface FormErrors {
  name?: string;
  email?: string;
  message?: string;
}

const ContactForm = () => {
  const [formData, setFormData] = useState<ContactFormData>({
    name: '',
    email: '',
    message: '',
  });

  const [errors, setErrors] = useState<FormErrors>({});
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [submitSuccess, setSubmitSuccess] = useState(false);
  const [submitError, setSubmitError] = useState<string | null>(null);

  // Validate form fields
  const validateForm = (): boolean => {
    const newErrors: FormErrors = {};

    // Name validation
    if (!formData.name.trim()) {
      newErrors.name = 'Имя обязательно';
    } else if (formData.name.trim().length < 2) {
      newErrors.name = 'Имя должно быть не менее 2 символов';
    } else if (formData.name.trim().length > 255) {
      newErrors.name = 'Имя не должно превышать 255 символов';
    }

    // Email validation
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!formData.email.trim()) {
      newErrors.email = 'Email обязателен';
    } else if (!emailRegex.test(formData.email)) {
      newErrors.email = 'Введите корректный email';
    } else if (formData.email.length > 255) {
      newErrors.email = 'Email не должен превышать 255 символов';
    }

    // Message validation
    if (!formData.message.trim()) {
      newErrors.message = 'Сообщение обязательно';
    } else if (formData.message.trim().length < 10) {
      newErrors.message = 'Сообщение должно быть не менее 10 символов';
    } else if (formData.message.length > 2000) {
      newErrors.message = 'Сообщение не должно превышать 2000 символов';
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  // Handle form submission
  const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setSubmitSuccess(false);
    setSubmitError(null);

    if (!validateForm()) {
      return;
    }

    setIsSubmitting(true);

    try {
      const response = await axios.post('/api/contact', formData);

      if (response.data.success) {
        setSubmitSuccess(true);
        setFormData({ name: '', email: '', message: '' });
        setErrors({});

        // Hide success message after 5 seconds
        setTimeout(() => {
          setSubmitSuccess(false);
        }, 5000);
      }
    } catch (error: any) {
      setSubmitError(
        error.response?.data?.message || 'Произошла ошибка при отправке сообщения. Попробуйте позже.'
      );
    } finally {
      setIsSubmitting(false);
    }
  };

  // Handle input changes
  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>
  ) => {
    const { name, value } = e.target;
    setFormData((prev) => ({
      ...prev,
      [name]: value,
    }));
    // Clear error for this field
    if (errors[name as keyof FormErrors]) {
      setErrors((prev) => ({
        ...prev,
        [name]: undefined,
      }));
    }
  };

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      className="bg-gradient-to-br from-gray-800 to-gray-900 rounded-2xl p-8 border-2 border-gray-700 shadow-2xl"
    >
      <form onSubmit={handleSubmit} className="space-y-6">
        {/* Name Field */}
        <div>
          <label
            htmlFor="name"
            className="block text-gray-300 font-semibold mb-2"
          >
            Имя *
          </label>
          <input
            type="text"
            id="name"
            name="name"
            value={formData.name}
            onChange={handleChange}
            className={`w-full px-4 py-3 bg-gray-700 border ${
              errors.name ? 'border-red-500' : 'border-gray-600'
            } rounded-lg text-white placeholder-gray-400 focus:outline-none focus:border-primary transition-colors`}
            placeholder="Ваше имя"
            disabled={isSubmitting}
          />
          {errors.name && (
            <p className="mt-2 text-sm text-red-500">{errors.name}</p>
          )}
        </div>

        {/* Email Field */}
        <div>
          <label
            htmlFor="email"
            className="block text-gray-300 font-semibold mb-2"
          >
            Email *
          </label>
          <input
            type="email"
            id="email"
            name="email"
            value={formData.email}
            onChange={handleChange}
            className={`w-full px-4 py-3 bg-gray-700 border ${
              errors.email ? 'border-red-500' : 'border-gray-600'
            } rounded-lg text-white placeholder-gray-400 focus:outline-none focus:border-primary transition-colors`}
            placeholder="your@email.com"
            disabled={isSubmitting}
          />
          {errors.email && (
            <p className="mt-2 text-sm text-red-500">{errors.email}</p>
          )}
        </div>

        {/* Message Field */}
        <div>
          <label
            htmlFor="message"
            className="block text-gray-300 font-semibold mb-2"
          >
            Сообщение *
          </label>
          <textarea
            id="message"
            name="message"
            value={formData.message}
            onChange={handleChange}
            rows={6}
            className={`w-full px-4 py-3 bg-gray-700 border ${
              errors.message ? 'border-red-500' : 'border-gray-600'
            } rounded-lg text-white placeholder-gray-400 focus:outline-none focus:border-primary transition-colors resize-none`}
            placeholder="Ваше сообщение..."
            disabled={isSubmitting}
          />
          <div className="flex justify-between items-center mt-2">
            {errors.message && (
              <p className="text-sm text-red-500">{errors.message}</p>
            )}
            <p
              className={`text-sm ${
                formData.message.length > 2000
                  ? 'text-red-500'
                  : 'text-gray-400'
              } ml-auto`}
            >
              {formData.message.length}/2000
            </p>
          </div>
        </div>

        {/* Success Message */}
        {submitSuccess && (
          <motion.div
            initial={{ opacity: 0, y: -10 }}
            animate={{ opacity: 1, y: 0 }}
            className="bg-green-600/20 border border-green-500 rounded-lg p-4"
          >
            <p className="text-green-400 flex items-center">
              <span className="text-2xl mr-2">✅</span>
              Спасибо! Ваше сообщение отправлено. Мы свяжемся с вами в ближайшее время.
            </p>
          </motion.div>
        )}

        {/* Error Message */}
        {submitError && (
          <motion.div
            initial={{ opacity: 0, y: -10 }}
            animate={{ opacity: 1, y: 0 }}
            className="bg-red-600/20 border border-red-500 rounded-lg p-4"
          >
            <p className="text-red-400 flex items-center">
              <span className="text-2xl mr-2">❌</span>
              {submitError}
            </p>
          </motion.div>
        )}

        {/* Submit Button */}
        <motion.button
          type="submit"
          whileHover={{ scale: 1.02 }}
          whileTap={{ scale: 0.98 }}
          disabled={isSubmitting}
          className={`w-full font-bold py-4 rounded-xl transition-all ${
            isSubmitting
              ? 'bg-gray-700 text-gray-400 cursor-not-allowed'
              : 'bg-gradient-to-r from-primary to-secondary text-white hover:shadow-lg hover:shadow-primary/50'
          }`}
        >
          {isSubmitting ? (
            <span className="flex items-center justify-center">
              <svg
                className="animate-spin h-5 w-5 mr-3"
                viewBox="0 0 24 24"
              >
                <circle
                  className="opacity-25"
                  cx="12"
                  cy="12"
                  r="10"
                  stroke="currentColor"
                  strokeWidth="4"
                  fill="none"
                />
                <path
                  className="opacity-75"
                  fill="currentColor"
                  d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                />
              </svg>
              Отправка...
            </span>
          ) : (
            '📧 Отправить сообщение'
          )}
        </motion.button>
      </form>
    </motion.div>
  );
};

export default ContactForm;

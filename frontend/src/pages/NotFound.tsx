import { Link } from 'react-router-dom';

const NotFound = () => {
  return (
    <div className="min-h-screen flex items-center justify-center bg-dark">
      <div className="text-center">
        <h1 className="text-9xl font-bold text-primary mb-4">404</h1>
        <h2 className="text-3xl font-semibold text-white mb-4">
          Страница не найдена
        </h2>
        <p className="text-gray-400 mb-8">
          К сожалению, запрашиваемая страница не существует
        </p>
        <Link
          to="/"
          className="bg-primary text-white font-semibold py-3 px-8 rounded-lg hover:bg-red-700 transition-colors inline-block"
        >
          Вернуться на главную
        </Link>
      </div>
    </div>
  );
};

export default NotFound;

import React,{useState} from 'react';
import axios from 'axios';
import './App.css';

const BASE_URL = 'http://localhost:8000';

const App: React.FC = () => {
  const [url, setUrl] = useState('');
  const [shortUrl, setShortUrl] = useState('');

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    try {
      const response = await axios.post(`${BASE_URL}/shorten`, {url});
      const { short_url } = response.data
      console.log(short_url);
      setShortUrl(`${BASE_URL}/${short_url}`);
    } catch (error) {
      console.log(error);
    }
  };

  const handleOpenUrl = () => {
    if (shortUrl) {
      window.open(shortUrl, '_blank');
    }
  };

  const handleCopyUrl = () => {
    if (shortUrl) {
      navigator.clipboard.writeText(shortUrl);
    }
  };

  return (
    <div className="container">
      <h1>Shorten URL</h1>
      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <input
            type="text"
            className="form-control"
            placeholder="Enter URL"
            value={url}
            onChange={(e) => setUrl(e.target.value)}
          />
        </div>
        <button type="submit" className="btn btn-primary">
          Shorten
        </button>
      </form>
      {shortUrl && (
        <div>
          <div>
            <p>Shortened URL : <a href={shortUrl} target="_blank" rel="noopener noreferrer">{shortUrl}</a></p>
          </div>
          <div>
            <button className="btn btn-primary" onClick={handleOpenUrl}>
              Open
            </button>
            <button className="btn btn-primary" onClick={handleCopyUrl}>
              Copy
            </button>
          </div>
        </div>)}
    </div>
  );
}

export default App;
import { useState, useEffect } from "react";
import axios from "axios";
import { useParams } from "react-router-dom";
import "../assets/styles/StockOverview.css";
import { Chart } from "react-google-charts";
import { LoremIpsum } from "react-lorem-ipsum";

function StockOverview() {
  const [data, setData] = useState(null);
  const [hist, setHist] = useState(null);

  let { id } = useParams();

  const createHistData = (hist) => {
    const len = Math.min(hist.values.length, hist.dates.length);

    let histData = [];

    histData.push(["Date", "Value"]);

    for (let i = 0; i < len; ++i) {
      histData.push([hist.dates[i], Math.round(hist.values[i] * 100) / 100]);
    }

    return histData;
  };

  useEffect(() => {
    const fetchData = async () => {
      try {
        const urlData = "http://localhost:2345/stock/info?symbol=" + id;
        const urlHist =
          "http://localhost:2345/stockHist?symbol=" + id + "&period=10y";
        const promisedData = axios.get(urlData);
        const promisedHist = axios.get(urlHist);
        const values = await Promise.all([promisedData, promisedHist]);
        const symbolData = values[0].data;
        const symbolHist = values[1].data;
        const histData = createHistData(symbolHist);
        console.log(symbolData);
        setHist(histData);
        setData(symbolData);
      } catch (error) {
        console.error(error.message);
      }
    };

    fetchData();
  }, []);

  const getImage = () => {
    if (data) {
      return <img style={{ marginRight: "5%" }} src={data.logoURL} />;
    }
  };

  const getCurrentOpen = () => {
    if (data) {
      return Math.round(data.currentOpen * 100) / 100;
    }
  };

  const getPredText = () => {
    if (!data) {
      return;
    }

    if (data.prediction === 1) {
      return (
        <div>
          Das Model hat mir einer Accuracy von{" "}
          {Math.round(data.accuracy * 100) / 100}% und einem Mean Squared Error
          von {Math.round(data.mse * 100) / 100} vorhergesagt, dass der Preis in 5 Tagen höher sein wird als er es heute ist. 
          Anhand von Vorhersage, wird Ihnen empfohlen, diese Aktie zu kaufen.
        </div>
      );
    }

    if (data.prediction === -1) {
      return (
        <div>
          Das Model hat mir einer Accuracy von{" "}
          {Math.round(data.accuracy * 100) / 100}% und einem Mean Squared Error
          von {Math.round(data.mse * 100) / 100} vorhergesagt, dass der Preis in 5 Tagen niedrieger sein wird als er es heute ist. 
          Anhand von Vorhersage, wird Ihnen empfohlen, diese Aktie zu verkaufen.
        </div>
      );
    }
  };

  const getPredImage = () => {
    if (!data) {
      return;
    }

    if (data.prediction === 1) {
      return (
        <img
          style={{ marginRight: "5%" }}
          src="data:image/png;base64,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"
        />
      );
    }

    if (data.prediction === 0) {
    }

    if (data.prediction === -1) {
      return (
        <img
          style={{ marginRight: "5%" }}
          src="data:image/png;base64,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"
        />
      );
    }
  };

  const getPreviousClose = () => {
    if (data) {
      return Math.round(data.previousClose * 100) / 100;
    }
  };

  const options = {
    title: "Historische Daten",
    curveType: "function",
    legend: { position: "bottom" },
  };

  return (
    <>
      <div className="parent">
        <div className="div1" style={{ marginBottom: "5%" }}>
          <div style={{ display: "flex" }}>
            {getImage()}
            <div>
              <h1>{data?.name}</h1>
              <div>
                Current open: {getCurrentOpen()}€ - Previous open:{" "}
                {getPreviousClose()}€
              </div>
              <div>
                Difference: {data?.diff}€ - Difference in Percent:{" "}
                {data?.diffPercent}%
              </div>
            </div>
          </div>
        </div>
        <div className="div2" style={{ marginBottom: "12%" }}>
          <Chart
            chartType="LineChart"
            width="100%"
            height="100%"
            data={hist}
            options={options}
          />
        </div>
        <div className="div3"> {getPredImage()} </div>
        <div className="div4" style={{ marginLeft: "15%" }}> {getPredText()} </div>
        <div className="div5">
          {" "}
          <LoremIpsum p={2} />{" "}
        </div>
        <div className="div6">
          {" "}
          {data?.description}{" "}
        </div>
      </div>
    </>
  );
}

export default StockOverview;

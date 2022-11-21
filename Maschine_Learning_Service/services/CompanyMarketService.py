from bs4 import BeautifulSoup
import requests
import bson

class CompanyMarketService:
    _baseUrl = str()

    def __init__(self) -> None:
        self._baseUrl = "https://companiesmarketcap.com/software/largest-software-companies-by-market-cap/?page="

    def GetCompanyMarketTicker(self) -> list:
        ticker_list = list()
        page = 1

        while True:

            if page >= 100:
                break

            url = self._baseUrl + str(page)
            response = requests.get(url)

            if response.status_code != 200:
                break
  
            soup = BeautifulSoup(response.text,'html.parser')

            trs = soup.find_all("tr")

            if len(trs) == 0:
                break
  
            del trs[0]

            if len(trs) == 0:
                break

            for tr in trs:
                companyNames = tr.find_all("div", {"class": "company-name"})
                if len(companyNames) == 0:
                    continue
    
                companySymbols = tr.find_all("div", {"class": "company-code"})
                if len(companySymbols) == 0:
                    continue

                data = {
                    "name": companyNames[0].text,
                    "symbol": companySymbols[0].text
                }

                ticker_list.append(data)

            page += 1

        return ticker_list

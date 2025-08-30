import akshare as ak


def get_stock_price(symbol: str):
    """
    获取个股最新价格
    :param symbol:
    :return:
    """
    stock_bid_ask_em_df = ak.stock_bid_ask_em(symbol=symbol)
    return stock_bid_ask_em_df.at[20, 'value']


def get_gainian_rank():
    """
    获取概念板块排行 top n
    :return:
    """
    stock_fund_flow_concept_df = ak.stock_fund_flow_concept(symbol="即时", max_page_num=1)
    return stock_fund_flow_concept_df.loc[:14]  # 50/page


def get_minute_price(symbol: str, period='1'):
    """
    获取概念板块分时行情
    :param symbol:
    :param period: 分钟粒度 1 5 10 ...
    :return:
    """
    stock_board_concept_hist_min_em_df = ak.stock_board_concept_hist_min_em(symbol=symbol, period=period)
    print(stock_board_concept_hist_min_em_df)


if __name__ == '__main__':
    # print(get_gainian_rank())
    print(get_minute_price('稀土永磁'))

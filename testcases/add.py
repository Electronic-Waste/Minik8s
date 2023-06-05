def main(params):
    x = params["x"]
    y = params["y"]
    result = {
        "res": x + 3 * y
    }
    print(result)
    return result

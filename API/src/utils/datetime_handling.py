from datetime import datetime


def convert_timestamp(date_string):
    format_code = "%Y-%m-%d"
    if date_string is None:
        return None 
    try:
        datetime_object = datetime.strptime(date_string, format_code)
        return datetime_object
    except ValueError:
        # This handles cases where the string doesn't match the format
        return None 
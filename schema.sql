CREATE TABLE verses (
        book_number NUMERIC NOT NULL, 
        chapter NUMERIC NOT NULL, 
        verse NUMERIC NOT NULL, 
        text TEXT NOT NULL DEFAULT '', 
        PRIMARY KEY (book_number, chapter, verse) );
CREATE TABLE info (
        name TEXT NOT NULL, 
        value TEXT NOT NULL, 
        PRIMARY KEY (name));
CREATE TABLE books (
        book_number NUMERIC NOT NULL, 
        short_name TEXT NOT NULL, 
        long_name TEXT NOT NULL, 
        book_color TEXT NOT NULL, 
        PRIMARY KEY (book_number));
CREATE TABLE books_all (
        book_number NUMERIC NOT NULL, 
        short_name TEXT NOT NULL, 
        long_name TEXT NOT NULL, 
        book_color TEXT NOT NULL, 
        is_present BOOLEAN NOT NULL, PRIMARY KEY (book_number));
CREATE TABLE stories (
        book_number NUMERIC NOT NULL, 
        chapter NUMERIC NOT NULL, 
        verse NUMERIC NOT NULL, 
        order_if_several NUMERIC NOT NULL DEFAULT 0, 
        title TEXT NOT NULL DEFAULT '', 
        PRIMARY KEY (book_number, chapter, verse, order_if_several));
CREATE TABLE introductions (
        book_number NUMERIC NOT NULL DEFAULT 0, 
        introduction TEXT NOT NULL DEFAULT '', 
        PRIMARY KEY (book_number));

class HomePresenter extends BasePresenter
  onCreate: () ->
    @searchBox = $('.js-search-box')

    @isSearchBoxCenterized = false
    @centerizeSearchBox()

    $('.js-input').on 'input change', (event) =>
      val = $(event.target).val()

      clearTimeout(@timer) if @timer?
      @request.abort() if @request?

      @timer = setTimeout(( () =>
        @timer = null
        @request.abort() if @request?

        if val.trim() == ""
          @centerizeSearchBox()
        else
          @moveSearchBoxToTop()

        keywords = val
        query = keywords.replace(/\+/gi, "").replace(/(\s)+/gi, "+")

        @request = $.get
          url: "/search?keywords=" + query
          success: (data) =>
            @request = null
            console.log JSON.parse(data.posts)
      ), 300)

  moveSearchBoxToTop: () ->
    return unless @isSearchBoxCenterized
    @searchBox.css('top', 0)
    @isSearchBoxCenterized = false

  centerizeSearchBox: () ->
    return if @isSearchBoxCenterized
    w = @searchBox.outerWidth()
    h = @searchBox.outerHeight()
    @searchBox.css 'top', (@searchBox.parent().height() - h) / 2
    @searchBox.css 'left', (@searchBox.parent().width() - w) / 2
    @isSearchBoxCenterized = true